let objs = [];
let sections = document.querySelectorAll('article.RichBlock-root');
sections.forEach(function (s) {
    let IMGobj = s.firstChild.firstChild.lastChild;
    let content = s.lastChild.firstChild.children;
    objs.push({
        URL: content[1].firstChild.firstChild.getAttribute('href'),
        IMGhref: IMGobj.getAttribute('src'),
        Tag: content[0].firstChild.innerText,
        Header:  content[1].firstChild.firstChild.firstChild.innerText
    })
}); 
return objs;