use std::io::{Write, BufWriter};

include!("dataset.rs");

// https://stackoverflow.com/a/51261570
fn char_windows<'a>(src: &str, win_size: usize) -> impl Iterator<Item = &str> {
  src.char_indices().flat_map(move |(from, _)| {
    src[from ..].char_indices()
      .nth(win_size - 1)
      .map(|(to, c)| {
        &src[from .. from + to + c.len_utf8()]
      })
  })
}

fn sliding_hashes(s: &str) -> Vec<u32> {
  // for t in char_windows(&s, 5) { println!("{}", t); }
  char_windows(&s, 5)
    .map(|s| s.chars().map(|c| c as u32)
      .fold(0u32, |x, y| x.wrapping_mul(997).wrapping_add(y)))
    .collect()
}

const HASH_N_DIMS: usize = 256;
fn charcode_hash(s: &str) -> [u16; HASH_N_DIMS] {
  let mut ret = [0; HASH_N_DIMS];
  for c in s.chars() {
    if c == '/' { continue; }
    ret[(c as usize) % HASH_N_DIMS] += 1;
  }
  ret
}

fn dist_manhattan(a: &[u16], b: &[u16]) -> u16 {
  debug_assert_eq!(a.len(), b.len());
  a.iter().zip(b.iter())
    .map(|(x, y)| (*x as i16 - *y as i16).abs())
    .sum::<i16>() as u16
}

// DSU
struct DSU {
  p: Vec<usize>,
}
impl DSU {
  fn new(n: usize) -> DSU {
    let mut p = Vec::with_capacity(n);
    for i in 0..n { p.push(i); }
    DSU { p }
  }

  fn root(&mut self, x: usize) -> usize {
    if self.p[x] != x {
      self.p[x] = self.root(self.p[x]);
    }
    self.p[x]
  }

  fn union(&mut self, x: usize, y: usize) {
    let x = self.root(x);
    let y = self.root(y);
    self.p[y] = x;
  }
}

fn main() {
  macro_rules! file {
    ($f:literal) => (concat!(
      "../../",
      // "../test/",
      $f))
  }
  const DATASET_PATH: &str = file!("1a-all.txt");
  const ACCEPT_PATH: &str = file!("20-accept.txt");
  const CURATE_PATH: &str = file!("21-curate.txt");
  const OUT_DUPS_PATH: &str = file!("2a-dups.txt");
  const OUT_DATA_PATH: &str = file!("2b-dedup.txt");

  let dataset = read_dataset(DATASET_PATH).unwrap();

  eprintln!("building hashes");
  let mut hashset = vec![];
  let mut map = std::collections::HashMap::new();
  for art in &dataset {
    let hashes = sliding_hashes(&art.content);
    for &h in &hashes {
      map.entry(h).or_insert(vec![]).push(art);
    }
    hashset.push(hashes);
    if art.id % 10000 == 0 { eprintln!("{}/{}", art.id, dataset.len()); }
  }

  eprintln!("checking for duplicates");
  let mut dsu = DSU::new(dataset.len());
  for (art, hashes) in dataset.as_slice().iter().zip(hashset.iter()) {
    let mut all: Vec<&Article> = vec![];
    hashes.iter().for_each(|h| all.extend(map.get(&h).unwrap()));
    all.sort_by_key(|art| art.id);
    all.dedup();
    // println!("{:?}", all.iter().map(|art| art.id).collect::<Vec<_>>());
    for other_art in all {
      if art == other_art { continue; }
      let dist = dist_manhattan(
        &charcode_hash(&art.content),
        &charcode_hash(&other_art.content),
      );
      if dist <= 20 {
        dsu.union(art.id, other_art.id);
        // println!("{}\n{}\n{}\n", dist, art.content, other_art.content);
      }
    }
    if art.id % 10000 == 0 { eprintln!("{}/{}", art.id, dataset.len()); }
  }

  eprintln!("marking curated articles");
  let mut contains = std::collections::HashMap::new();
  let mut is_curated = vec![false; dataset.len()];
  for art in &dataset {
    for s in art.content.split('/') {
      contains.entry(s.to_string()).or_insert(vec![]).push(art.id);
    }
  }
  let curated = (|| -> UniResult<_> {
    Ok(BufReader::new(std::fs::File::open(CURATE_PATH)?)
      .lines().flat_map(|s|
        Some(s.ok()?.split(|c|
            c == '，' || c == '。' || c == '、' ||
            c == '？' || c == '！' || c == '：')
          .flat_map(|s| if s.is_empty() { None } else { Some(s.to_string()) })
          .collect::<Vec<_>>()))
      .collect::<Vec<_>>())
  })().unwrap_or(vec![]);
  for line in curated {
    let mut candidates = std::collections::HashSet::<usize>::new();
    for s in &line {
      match contains.get(&s.to_string()) {
        Some(list) => { candidates.extend(list); },
        None => (),
      }
    }
    let joined = line.join("/");
    let mut best = (usize::MAX, usize::MAX);
    for id in candidates {
      if dataset[id].content.contains(&joined) {
        let cur_len = dataset[id].content.len();
        if cur_len < best.1 {
          best = (id, cur_len);
        }
      }
    }
    if best.1 < usize::MAX {
      // eprintln!("{} {} {}", joined, dsu.root(best.0), dataset[best.0].content);
      is_curated[dsu.root(best.0)] = true;
    }
  }

  eprintln!("printing/accepting duplicates");
  let accept = (|| -> UniResult<_> {
    Ok(BufReader::new(std::fs::File::open(ACCEPT_PATH)?)
      .lines().flat_map(|s| {
        let s = s.ok()?;
        if s.is_empty() || s.starts_with('#') { None } else { Some(s) }
      })
      .collect::<std::collections::HashSet<_>>())
  })().unwrap_or(std::collections::HashSet::new());
  let mut f_dups =
    BufWriter::new(std::fs::File::create(OUT_DUPS_PATH).unwrap());
  let mut f_data =
    BufWriter::new(std::fs::File::create(OUT_DATA_PATH).unwrap());

  let mut groups = vec![vec![]; dataset.len()];
  for i in 0..dataset.len() { groups[dsu.root(i)].push(i); }
  for i in 0..dataset.len() {
    if i % 10000 == 0 { eprintln!("{}/{}", i, dataset.len()); }
    if groups[i].is_empty() { continue; }
    let mut g = groups[i].iter()
      .map(|&id| (id, &dataset[id].content))
      .collect::<Vec<_>>();
    g.sort_by_key(|a| a.1);
    g.dedup_by_key(|a| a.1);
    let gacc = g.iter().filter(|s| accept.contains(&**s.1)).collect::<Vec<_>>();
    let marker = |j| if j == 0 {
      if is_curated[i] { '!' } else { '*' }
    } else { ' ' };
    if !gacc.is_empty() {
      for (i, a) in gacc.iter().enumerate() {
        writeln!(f_data, "{}\t{}", marker(i), dataset[a.0]).unwrap();
      }
    } else if g.len() > 1 {
      for (i, a) in g.iter().enumerate() {
        writeln!(f_dups, "{}", a.1).unwrap();
        writeln!(f_data, "{}\t{}", marker(i), dataset[a.0]).unwrap();
      }
      writeln!(f_dups).unwrap();
    } else {
      writeln!(f_data, "{}\t{}", marker(0), dataset[i]).unwrap();
    }
  }

  eprintln!("done!");
}
