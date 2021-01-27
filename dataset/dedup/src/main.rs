use std::io::{BufRead, BufReader};
use kdtree::KdTree;

type UniResult<T> = Result<T, Box<dyn std::error::Error>>;

#[derive(Debug, Clone)]
struct Article {
  dynasty: String,
  author: String,
  title: String,
  content: String,
}

fn read_dataset(path: &str) -> UniResult<Vec<Article>> {
  let lines = BufReader::new(std::fs::File::open(path)?).lines();
  let mut articles = vec![];
  for line in lines {
    let line = line?;
    let fields = line.split('\t').collect::<Vec<_>>();
    articles.push(Article {
      dynasty: fields[0].to_string(),
      author: fields[1].to_string(),
      title: fields[2].to_string(),
      content: fields[3].to_string(),
    });
  }
  Ok(articles)
}

const HASH_N_DIMS: usize = 64;

fn calc_hash(s: &str) -> [u16; HASH_N_DIMS] {
  let mut ret = [0; HASH_N_DIMS];
  for c in s.chars() {
    if c == '/' { continue; }
    ret[(c as usize) % HASH_N_DIMS] += 1;
  }
  ret
}

// Make this a procedural macro
macro_rules! array_convert {
  ($a: expr, $t: ty, $n: expr) => {
    [
      $a[ 0] as $t, $a[ 1] as $t, $a[ 2] as $t, $a[ 3] as $t,
      $a[ 4] as $t, $a[ 5] as $t, $a[ 6] as $t, $a[ 7] as $t,
      $a[ 8] as $t, $a[ 9] as $t, $a[10] as $t, $a[11] as $t,
      $a[12] as $t, $a[13] as $t, $a[14] as $t, $a[15] as $t,
      $a[16] as $t, $a[17] as $t, $a[18] as $t, $a[19] as $t,
      $a[20] as $t, $a[21] as $t, $a[22] as $t, $a[23] as $t,
      $a[24] as $t, $a[25] as $t, $a[26] as $t, $a[27] as $t,
      $a[28] as $t, $a[29] as $t, $a[30] as $t, $a[31] as $t,
      $a[32] as $t, $a[33] as $t, $a[34] as $t, $a[35] as $t,
      $a[36] as $t, $a[37] as $t, $a[38] as $t, $a[39] as $t,
      $a[40] as $t, $a[41] as $t, $a[42] as $t, $a[43] as $t,
      $a[44] as $t, $a[45] as $t, $a[46] as $t, $a[47] as $t,
      $a[48] as $t, $a[49] as $t, $a[50] as $t, $a[51] as $t,
      $a[52] as $t, $a[53] as $t, $a[54] as $t, $a[55] as $t,
      $a[56] as $t, $a[57] as $t, $a[58] as $t, $a[59] as $t,
      $a[60] as $t, $a[61] as $t, $a[62] as $t, $a[63] as $t,
    ]
  };
}

fn dist_manhattan(a: &[f32], b: &[f32]) -> f32 {
  debug_assert_eq!(a.len(), b.len());
  a.iter().zip(b.iter())
    .map(|(x, y)| (*x - *y).abs())
    .fold(0.0, |x, y| { x + y })
}

fn main() {
  let dataset = read_dataset("../../all.txt").unwrap();
  // let dataset = dataset[..10].to_vec();
  // println!("{:?}", dataset);

  println!("building tree");
  let mut tree = KdTree::new(HASH_N_DIMS);
  let mut hashes = vec![];
  for (i, article) in dataset.iter().enumerate() {
    let hash = calc_hash(&article.content);
    let hashf = array_convert!(hash, f32, HASH_N_DIMS);
    hashes.push(hashf);
    tree.add(hashf, i).unwrap();
    if i % 10000 == 0 { println!("{}/{}", i, dataset.len()); }
  }

  println!("checking for duplicates");
  let mut i = 0;
  for (article, hash) in dataset.iter().zip(hashes.iter()) {
    let nghbrs = tree.within(hash, 12.0, &dist_manhattan).unwrap();
    let nghbrs = nghbrs.iter().filter(|x| x.0 > 0.0).collect::<Vec<_>>();
    if nghbrs.len() > 0 {
      println!("{}", article.content);
      for (dist, idx) in nghbrs {
        println!("{} {:?}", dist, dataset[**idx as usize].content);
      }
      println!("");
    }
    i += 1;
    if i % 100 == 0 { println!("{}/{}", i, dataset.len()); }
  }
}
