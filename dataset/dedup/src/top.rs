include!("dataset.rs");

fn main() {
  // const DATASET_PATH: &str = "test.txt";
  const DATASET_PATH: &str = "../../all.txt";

  let dataset = read_dataset(DATASET_PATH).unwrap();

  eprintln!("counting characters");
  let mut map = std::collections::HashMap::new();
  for art in &dataset {
    for c in art.content.chars() {
      if c == '/' { continue; }
      *map.entry(c).or_insert(0) += 1;
    }
    if art.id % 10000 == 0 { eprintln!("{}/{}", art.id, dataset.len()); }
  }

  let mut vec = map.iter().collect::<Vec<_>>();
  vec.sort_by_key(|a| -a.1);
  vec.truncate(500);
/*
  println!("{}", vec.iter()
    .map(|(ch, cnt)| format!("{} {}", ch, cnt))
    .collect::<Vec<_>>()
    .join("\n"));
*/
  println!("{}", vec.iter()
    .map(|(ch, _cnt)| ch)
    .collect::<Vec<_>>()
    .as_slice().chunks(25)
    .map(|s| s.iter().map(|c| **c).collect::<String>())
    .collect::<Vec<_>>()
    .join("\n"));

  eprintln!("done!");
}
