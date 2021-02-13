export interface BlogBrief {
    id: string;
    title: string;
    cover: string;
    rate: number;
    description: string;
    createTime: string;
    author: string;
}
export interface User {
    id: string;
    name: string;
    email: string;
    description: string;
    createTime: string;
    site?: string;
    avatar: string;
}