import { httpProtocol, serverURL } from "@/config/defaults"

const url = httpProtocol+serverURL+"/uploads"

type ImageProps = {
    src: string
    width: number
    quality?: number
}

export default function ImageLoader({ src, width, quality}: ImageProps) {
    return `${url}/${src}?w=${width}&q=${quality || 100}`
}