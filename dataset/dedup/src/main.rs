use std::io::{BufRead, BufReader, Write, BufWriter};

type UniResult<T> = Result<T, Box<dyn std::error::Error>>;

#[derive(Debug, Clone)]
struct Article {
  id: usize,
  dynasty: String,
  author: String,
  title: String,
  content: String,
}
impl PartialEq for Article {
  fn eq(&self, other: &Self) -> bool { self.id == other.id }
}
impl Eq for Article { }
impl std::fmt::Display for Article {
  fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
    write!(f, "{}\t{}\t{}\t{}",
      self.dynasty, self.author, self.title, self.content)
  }
}

fn read_dataset(path: &str) -> UniResult<Vec<Article>> {
  let lines = BufReader::new(std::fs::File::open(path)?).lines();
  let mut articles = vec![];
  for line in lines {
    let line = line?;
    let fields = line.split('\t').collect::<Vec<_>>();
    articles.push(Article {
      id: articles.len(),
      dynasty: fields[0].to_string(),
      author: fields[1].to_string(),
      title: fields[2].to_string(),
      content: fields[3].to_string(),
    });
  }
  Ok(articles)
}

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
  // const DATASET_PATH: &str = "test.txt";
  // const ACCEPT_PATH: &str = "test_accept.txt";
  const DATASET_PATH: &str = "../../all.txt";
  const ACCEPT_PATH: &str = "../../all_accept.txt";
  const OUT_DUPS_PATH: &str = "dups.txt";
  const OUT_DATA_PATH: &str = "dedup.txt";

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
    if !gacc.is_empty() {
      for a in gacc { writeln!(f_data, "{}", dataset[a.0]).unwrap(); }
    } else if g.len() > 1 {
      for a in g {
        writeln!(f_dups, "{}", a.1).unwrap();
        writeln!(f_data, "{}", dataset[a.0]).unwrap();
      }
      writeln!(f_dups).unwrap();
    } else {
      writeln!(f_data, "{}", dataset[i]).unwrap();
    }
  }

  eprintln!("done!");
}
