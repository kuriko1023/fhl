use std::io::{BufRead, BufReader};

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
impl PartialOrd for Article {
  fn partial_cmp(&self, other: &Self) -> Option<std::cmp::Ordering> {
    self.id.partial_cmp(&other.id)
  }
}
impl Ord for Article {
  fn cmp(&self, other: &Self) -> std::cmp::Ordering { self.id.cmp(&other.id) }
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
    .map(|s| s.chars().map(|c| c as u32).fold(0, |x, y| x * 997 + y))
    .collect()
}

/*
const HASH_WIN: usize = 5;
const HASH_BASE: u32 = 997;
fn sliding_hashes(s: &str) -> Vec<u32> {
  let mut window = [0u32; HASH_WIN];
  let mut hash = 0;
  let mut hashes = vec![];
  for (i, ch) in s.chars().enumerate() {
    hash -= window[0] * HASH_BASE.pow(HASH_WIN as u32);
    window.rotate_left(1);
    window[HASH_WIN - 1] = ch as u32;
    hash = hash * HASH_BASE + ch as u32;
    if i >= HASH_WIN - 1 { hashes.push(hash); }
  }
  hashes
}
*/

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
  // let dataset = read_dataset("../../all.txt").unwrap();
  // let dataset = dataset[..10].to_vec();
  // println!("{:?}", dataset);
  let dataset = read_dataset("test.txt").unwrap();

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
    all.sort();
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

  eprintln!("printing duplicates");
  let mut groups = vec![vec![]; dataset.len()];
  for i in 0..dataset.len() { groups[dsu.root(i)].push(i); }
  for i in 0..dataset.len() {
    let mut g = groups[i].iter()
      .map(|&id| &dataset[id].content)
      .collect::<Vec<_>>();
    g.sort();
    g.dedup();
    if g.len() > 1 {
      for s in g { println!("{}", s); }
      println!();
    }
    if i % 10000 == 0 { eprintln!("{}/{}", i, dataset.len()); }
  }

  eprintln!("done!");
}
