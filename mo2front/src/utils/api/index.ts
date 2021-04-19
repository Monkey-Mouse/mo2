import {
    Blog,
    BlogBrief,
    BlogUpsert,
    Category,
    Count,
    ImgToken,
    SubComment,
    User,
    UserListData,
    Comment,
    Notification
} from "@/models";
import axios from "axios";
import * as qiniu from 'qiniu-js';
import { GetErrorMsg } from "..";
export * from './like'

function onlyUnique(value, index, self) {
    return self.indexOf(value) === index;
}

export async function GetUserInfoAsync() {
    let re = await axios.get<User>('/api/logs');
    return re.data
}
export async function RegisterAsync(userInfo: { email: string; password: string; userName: string }) {
    return (await axios.post<User>('/api/accounts', userInfo)).data;
}
export async function LoginAsync(userInfo: { userNameOrEmail: string; password: string }) {
    return (await axios.post<User>('/api/accounts/login', userInfo)).data;
}
export async function GetUserData(uid: string): Promise<User> {
    let re = await axios.get<User>('/api/accounts/detail/' + uid);
    return re.data[0]
}
export async function GetUserDatas(uids: string[]): Promise<UserListData[]> {
    if (uids.length === 0) {
        return [];
    }
    let re = await axios.get<UserListData[]>('/api/accounts/listBrief?id=' + uids.filter(onlyUnique).join('&id='));
    return re.data ?? []
}
export async function GetUploadToken(fname: string) {
    return (await axios.get<ImgToken>('/api/img/' + fname)).data
}

export function ParseQuery(query: { [key: string]: any }) {
    let queryStr = '?';
    const queryList: string[] = [];
    for (const key in query) {
        const element = query[key];
        queryList.push(`${key}=${element}`)
    }
    queryStr = queryStr + queryList.join('&');
    return queryStr
}
export const GetArticles = async (query: { page: number; pageSize: number; draft: boolean; search?: string }) => {
    return (await axios.get<BlogBrief[]>('/api/blogs/query' + ParseQuery(query))).data
}
export async function UpsertBlog(query: { draft: boolean }, blog: BlogUpsert) {
    if (!blog.categories || blog.categories.length === 0) {
        blog.categories = []
    }
    return (await axios.post<Blog>('/api/blogs/publish' + ParseQuery(query), blog)).data
}
export function UpSertBlogSync(query: { draft: boolean }, blog: BlogUpsert) {

    navigator.sendBeacon("/api/blogs/publish" + ParseQuery(query), JSON.stringify(blog))
}
export async function GetArticle(query: { id: string; draft: boolean }) {
    return (await axios.get<Blog>('/api/blogs/find/id' + ParseQuery(query))).data
}
export const GetOwnArticles = async (query: { page: number; pageSize: number; draft: boolean, deleted?: boolean }) => {
    return (await axios.get<BlogBrief[]>('/api/blogs/find/own' + ParseQuery(query))).data
}

export const GetUserArticles = async (query: { page: number; pageSize: number; draft: boolean; id: string }) => {
    return (await axios.get<BlogBrief[]>('/api/blogs/find/userId' + ParseQuery(query))).data
}
export async function DeleteArticle(id: string, query: { draft: boolean }) {
    (await axios.delete('/api/blogs/' + id + ParseQuery(query)))
}
export async function Logout() {
    (await axios.post('/api/accounts/logout'));
}
export async function UpdateUserInfo(info: User) {
    return (await axios.put<User>('/api/accounts', info)).data;
}
export async function UploadMD(md: File) {
    let form = new FormData();
    form.append('upload[]', md)
    return (await axios.post<Blog>('/api/file', form)).data;
}
export async function GetCategories(id: string) {
    return (await axios.get<Category[]>('/api/relation/category/sub/' + id)).data ?? []
}
export async function DeleteCategories(ids: string[]) {
    await axios.delete('/api/directories/category', { data: ids });
}
export async function UpsertCate(cate: Category) {
    return await (await axios.post<Category>("/api/blogs/category", cate)).data
}

export async function GetCateBlogs(id: string) {
    return (await axios.get<Blog[]>('/api/relation/blogs/category/' + id)).data ?? []
}

export async function GetCates(id: string) {
    return (await axios.get<Category[]>('/api/relation/category/user/' + id)).data ?? []
}

export async function GetComments(id: string, query: { page: number; pagesize: number }) {
    return (await axios.get<Comment[]>('/api/comment/' + id + ParseQuery(query))).data ?? []
}
export async function GetCommentNum(id: string) {
    return (await axios.get<Count>('/api/commentcount/' + id)).data
}
export async function UpsertComment(c: Comment) {
    return (await axios.post<Comment>('/api/comment', c)).data
}
export async function UpsertSubComment(id: string, c: SubComment) {
    return (await axios.post<SubComment>('/api/comment/' + id, c)).data
}
export async function GetNotificationNums() {
    return (await axios.get<{ num: number }>("/api/notification/num")).data
}
export async function GetNotifications(query: { page: number, pagesize: number }) {
    return (await axios.get<Notification[]>("/api/notification" + ParseQuery(query))).data
}
export async function RecycleBlog(id: string, query: { draft: boolean }) {
    axios.put('/api/blogs/recycle/' + id + ParseQuery(query))
}
export async function RestoreBlog(id: string, query: { draft: boolean }) {
    axios.put('/api/blogs/restore/' + id + ParseQuery(query))
}
