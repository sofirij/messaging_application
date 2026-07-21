const url = `http://${process.env.NEXT_PUBLIC_API_URL}/uploads`

type ImageProps = {
    src: string
    width: number
    quality?: number
}

export default function ImageLoader({ src, width, quality}: ImageProps) {
    return `${url}/${src}?w=${width}&q=${quality || 100}`
}