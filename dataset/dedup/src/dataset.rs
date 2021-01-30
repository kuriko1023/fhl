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


