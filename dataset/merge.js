const fs = require('fs');

// awk 'BEGIN { FS = "\t" } { print $1 }' < all.txt | sort -u
const normalizeDynasty = function (s) {
  if (s.match(/[近现当]/g)) return '近现代';
  if (s.match(/五代/g)) return '五代十国';
  return s.trim().replace(/[朝代]/, '');
};

const specCharList = [];
const specChar = function (s) {
  for (const ch of s) {
    const c = ch.charCodeAt(0);
    if (!(
      (c >= 0x4E00 && c <= 0x9FEF) ||
      (c >= 0x3400 && c <= 0x4DB5) ||
      (c >= 0x20000 && c <= 0x2A6D6) ||
      (c >= 0x2A700 && c <= 0x2B734) ||
      (c >= 0x2B740 && c <= 0x2B81D) ||
      (c >= 0x2B820 && c <= 0x2CEA1) ||
      (c >= 0x2CEB0 && c <= 0x2EBE0) ||
      (c >= 0x30000 && c <= 0x3134A)
    )) {
      if (c >= 0xD800 && c <= 0xDB7F) {
        // High Surrogates
      } else if (c >= 0xE000 && c <= 0xF8FF) {
        // Private Use Area
      } else {
        if (specCharList.indexOf(c) === -1) {
          specCharList.push(c);
          console.log(ch + '\t' + c.toString(16) + '\t' + s);
        }
      }
      return true;
    }
  }
  return false;
};

const process = function (s) {
  const a = s
    .replace(/\(.+?\)/g, '')
    .replace(/（.+?）/g, '')
    .replace(/[、《》「」『』“”　 ]/g, '')
    .replace(/<.+?>/g, '')
    .replace(/[⓪①②③④⑤⑥⑦⑧⑨⑩⑪⑫⑬⑭⑮⑯⑰⑱⑲⑳⑴⑵⑶⑷⑸⑹⑺⑻⑼⑽⑾⑿⒀⒁⒂⒃⒄⒅⒆⒇]/g, '')
    .split(/[。，：；？！.,:;?!\n]+/)
    .filter(e => e !== '');
  if (a.some(e => specChar(e))) return null;
  return a;
};

const all = [];

console.log('chinese-poetry/chinese-poetry');
for (const f of fs.readdirSync('./10-sources/1/wudai/huajianji')) {
  if (f.indexOf('juan') === -1) continue;
  for (const item of
    JSON.parse(fs.readFileSync('./10-sources/1/wudai/huajianji/' + f)))
  {
    if (item.author === '尹鹗 六首') item.author = '尹鹗';
    if (item.author === '欧陽炯') item.author = '欧阳炯';
    if (item.title === '河传·曲槛') continue;
    all.push([
      (item.author === '温庭筠' || item.author === '皇甫松') ?
        '唐' : '五代',
      item.author,
      item.title,
      item.paragraphs.join('\n').replace('B*', '圭')
    ]);
  }
}

console.log('javayhu/poetry');
for (const f of fs.readdirSync('./10-sources/2/data/poetry')) {
  const item = JSON.parse(fs.readFileSync('./10-sources/2/data/poetry/' + f));
  all.push([
    item.dynasty,
    item.poet.name,
    item.name,
    item.content,
  ]);
}

console.log('yxcs/poems-db');
for (let i = 1; i <= 4; i++) {
  const lines = fs.readFileSync('./10-sources/3/poems' + i + '.json')
    .toString().trim().split('\n');
  for (const s of lines) {
    const item = JSON.parse(s);
    if (!item.content) continue;
    if (item.author === '余光中') continue;
    all.push([
      item.dynasty,
      item.author,
      item.name,
      item.content.join('\n'),
    ]);
  }
}

console.log('snowtraces/poetry-source');
for (const name of '诗词曲') {
  const path = './10-sources/4/source/' + name;
  for (const d1 of fs.readdirSync(path)) if (d1 !== 'README.md') {
    for (const f of fs.readdirSync(path + '/' + d1)) {
      const file = path + '/' + d1 + '/' + f;
      const arr = JSON.parse(fs.readFileSync(file));
      for (const item of arr) if (item.content)
        all.push([
          item.dynasty,
          item.authorName,
          item.title,
          item.content.join('\n'),
        ]);
    }
  }
}
for (const item of
  JSON.parse(fs.readFileSync('./10-sources/4/source/其他/诗经.json')))
{
  all.push([
    '先秦',
    '佚名',
    item.title,
    item.content.join('\n'),
  ]);
}
for (const f of fs.readdirSync('./10-sources/4/全唐诗/ZZU_JSON_chs')) {
  for (const item of
    JSON.parse(fs.readFileSync('./10-sources/4/全唐诗/ZZU_JSON_chs/' + f)))
  {
    all.push([
      '唐',
      item.author,
      item.title,
      item.content.join('\n'),
    ]);
  }
}

for (const line of fs.readFileSync('./11-extra.txt').toString().split('\n')) {
  if (line) all.push(line.split('\t'));
}

let collCount = 0;
let errCount = 0;
const fdAll = fs.openSync('./1a-all.txt', 'w');
const fdErr = fs.openSync('./1b-err.txt', 'w');
all.forEach((x) => {
  const dynasty = normalizeDynasty(x[0]);
  const pars = process(x[3]);
  if (!x[0] || !x[1] || !x[2] || pars === null) {
    fs.writeSync(fdErr,
      x[0] + '\t' + x[1] + '\t' + x[2] + '\n' + x[3] + '\n\n');
    errCount++;
    return;
  }
  fs.writeSync(fdAll,
    dynasty + '\t' + x[1].trim() + '\t' + x[2].trim() + '\t' +
    pars.join('/') + '\n');
  collCount++;
});
fs.closeSync(fdAll);
fs.closeSync(fdErr);

console.log('collected ' + collCount);
console.log('filtered  ' + errCount);
