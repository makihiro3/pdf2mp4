// input type="file"要素からファイルのバイト列を取り出す
const readFileAsync = file =>
    new Promise((resolve, reject) => {
        const reader = new FileReader()
        reader.onload = () => resolve(reader.result)
        reader.onerror = () => reject(reader.error)
        reader.readAsArrayBuffer(file)
    })

const downloadLink = (blob, filename) => {
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.style = 'display: none'
    a.href = url
    a.download = filename
    document.body.appendChild(a)
    a.click()
    a.remove()
}

const convert = async (file, size) => {
    const res = await fetch(`/convert.cgi?size=${size}`,{
        method: "POST",
        mode: "same-origin",
        cache: "no-store",
        credentials: "same-origin",
        redirect: "error",
        headers: {
            "Content-Type": "application/pdf",
        },
        body: file,
    })
    if(res.status != 200) {
        console.error("error in backend", res)
        return null
    }
    return res.blob()
}

const handler = async event =>{
    const size = document.getElementById("size").value
    const input = document.getElementById("input")
    if(!(input.files instanceof FileList)){
        console.error("input is not input[file] node")
        return
    }
    if(input.files.length == 0) {
        console.warn("require file")
        return
    }
    for(var f of input.files){
        let name = f.name
        if(name.endsWith(".pdf")){
            name = name.substring(0, name.length-4)
        }
        name += `.${size}.mp4`
        const content = await readFileAsync(f)
        const blob = await convert(content, size)
        if (blob !== null) {
            downloadLink(blob, name)
        }
    }
}

const main = async event => {
    document.getElementById("form").addEventListener("submit", handler)
}
document.addEventListener("DOMContentLoaded", main)
