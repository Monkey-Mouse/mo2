import axios from "axios";

export async function GetBlogLikeNum(id: string) {
    return (await axios.get<{ num: number }>('/api/like/num/blog/' + id)).data;
}

export async function ToggleLikeBlog(id: string) {
    await axios.post('/api/like/blog', { "id": id });
}

export async function IsBlogLiked(id: string) {
    return (await axios.get<{ liked: boolean }>('/api/like/ext/blog/' + id)).data;
}
