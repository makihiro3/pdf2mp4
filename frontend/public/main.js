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
    a.click()
    a.remove()
}

const convert = async (file, size, interval) => {
    const res = await fetch(`/convert.cgi?size=${size}&interval=${interval}`,{
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
        if (res.status == 400) {
            log("client side error")
        } else {
            log("server side error")
        }
        return null
    }
    return res.blob()
}

const handler = async event =>{
    const size = document.getElementById("size").value
    const interval = document.getElementById("interval").value
    const input = document.getElementById("input")
    if(!(input.files instanceof FileList)){
        log("input is not input[file] node")
        return
    }
    if(input.files.length == 0) {
        log("ファイルを指定してください")
        return
    }
    for(const f of input.files){
        let name = f.name
        log(`変換開始 name:"${name}" size:${size} interval:${interval}`)
        if(name.endsWith(".pdf")){
            name = name.substring(0, name.length-4)
        }
        name += `.t${interval}.r${size}.mp4`
        const content = await readFileAsync(f)
        const blob = await convert(content, size, interval)
        if (blob !== null) {
            downloadLink(blob, name)
            log(`変換終了 name:"${name}"`)
        }
    }
}

const log = (message) => {
    const now = new Date()
    console.log(message)
    document.getElementById("log").value += `${now.toLocaleString()} ${message}\n`
}

const main = async event => {
    document.getElementById("form").addEventListener("submit", handler)
}
document.addEventListener("DOMContentLoaded", main)
