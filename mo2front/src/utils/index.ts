import { User, ApiError, ImgToken, BlogBrief, BlogUpsert, Blog, UserListData, Category } from '@/models/index'
import axios, { AxiosError } from 'axios';
import * as qiniu from 'qiniu-js';
import router from '../router'

export function randomProperty(obj: any) {
    const keys = Object.keys(obj);
    return obj[keys[keys.length * Math.random() << 0]];
}

export function Copy<T>(mainObject: T) {
    const objectCopy = {}; // objectCopy will store a copy of the mainObject
    let key;
    for (key in mainObject) {
        objectCopy[key] = mainObject[key]; // copies each property to the objectCopy object
    }
    return objectCopy as T;
}
export async function GetUserData(uid: string): Promise<User> {
    let re = await axios.get<User>('/api/accounts/detail/' + uid);
    return re.data[0]
}
function onlyUnique(value, index, self) {
    return self.indexOf(value) === index;
}
export async function GetUserDatas(uids: string[]): Promise<UserListData[]> {
    let re = await axios.get<UserListData[]>('/api/accounts/listBrief?id=' + uids.filter(onlyUnique).join('&id='));
    return re.data
}

export function GetInitials(name: string) {
    let rgx = new RegExp(/(\p{L}{1})\p{L}+/, 'gu');

    let initials = [...name.matchAll(rgx)] || [];

    return (
        (initials.shift()?.[1] || '') + (initials.pop()?.[1] || '')
    ).toUpperCase();
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
export function GetErrorMsg(apiError: any) {
    const err = (apiError as AxiosError<ApiError>);
    try {
        if (err.response.status === 404) {
            router.push('/404')
        }
        return err.response.data.reason
    } catch (error) {
        return 'Unknown Error'
    }
}
export async function GetUploadToken(fname: string) {
    return (await axios.get<ImgToken>('/api/img/' + fname)).data
}
export const UploadImgToQiniu = async (
    blobs: File[],
    callback: (imgprop: { src: string }) => void
) => {
    const promises: Promise<void>[] = []
    for (let index = 0; index < blobs.length; index++) {
        const element = blobs[index];
        const promise = new Promise<void>((resolve, reject) => {
            GetUploadToken(element.name).then(val => {
                let ob = qiniu.upload(element, val.file_key, val.token);
                ob.subscribe(null, (err) => {
                    reject(err)
                }, res => {
                    callback({ src: '//cdn.mo2.leezeeyee.com/' + res.key })
                    resolve();
                })
            })
        })
        promises.push(promise)
    }
    await Promise.all(promises)

}
export var globaldic: any = {};
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
export const GetArticles = async (query: { page: number; pageSize: number; draft: boolean }) => {
    return (await axios.get<BlogBrief[]>('/api/blogs/query' + ParseQuery(query))).data
}
export async function UpsertBlog(query: { draft: boolean }, blog: BlogUpsert) {
    return (await axios.post<Blog>('/api/blogs/publish' + ParseQuery(query), blog)).data
}
export function UpSertBlogSync(query: { draft: boolean }, blog: BlogUpsert) {

    navigator.sendBeacon("/api/blogs/publish" + ParseQuery(query), JSON.stringify(blog))
}
export async function GetArticle(query: { id: string; draft: boolean }) {
    return (await axios.get<Blog>('/api/blogs/find/id' + ParseQuery(query))).data
}
export const GetOwnArticles = async (query: { page: number; pageSize: number; draft: boolean }) => {
    return (await axios.get<BlogBrief[]>('/api/blogs/find/own' + ParseQuery(query))).data
}

export const GetUserArticles = async (query: { page: number; pageSize: number; draft: boolean; id: string }) => {
    return (await axios.get<BlogBrief[]>('/api/blogs/find/userId' + ParseQuery(query))).data
}
export function ReachedBottom(): boolean {
    return (window.innerHeight + window.pageYOffset) >= document.body.offsetHeight;
}
export interface BlogAutoLoader {
    blogs: BlogBrief[];
    loading: boolean;
    firstloading: boolean;
    page: number;
    pagesize: number;
    nomore: boolean;
    ReachedButtom: () => void;
}
export function ElmReachedButtom(elm: BlogAutoLoader, getArticles: (query: { page: number; pageSize: number }) => Promise<BlogBrief[]>) {
    if (elm.loading === false && !elm.nomore) {
        elm.loading = true;
        getArticles({
            page: elm.page++,
            pageSize: elm.pagesize,
        }).then((val) => {
            try {
                AddMore(elm, val);
            } catch (error) {
                elm.loading = false;
            }
        });
    }
}
export function AddMore(elm: BlogAutoLoader, val: BlogBrief[]) {
    if (!val || val.length < elm.pagesize) {
        elm.nomore = true;
    }
    if (!val) {
        elm.loading = false;
        return
    }
    for (let index = 0; index < val.length; index++) {
        const element = val[index];
        elm.blogs.push(element);
    }
    elm.loading = false;
}
export async function DeleteArticle(id: string, query: { draft: boolean }) {
    (await axios.delete('/api/blogs/' + id + ParseQuery(query)))
}
export async function Logout() {
    (await axios.post('/api/accounts/logout'));
}
export const AdminRole = "GeneralAdmin"
export const UserRole = "OrdinaryUser"
export const AnonymousRole = "Anonymous"
export function timeout(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}
export async function UpdateUserInfo(info: User) {
    return (await axios.put<User>('/api/accounts', info)).data;
}
export async function UploadMD(md: File) {
    let form = new FormData();
    form.append('upload[]', md)
    return (await axios.post<Blog>('/api/file', form)).data;
}
export async function addQuery(that: Vue, key: string, val: string | string[]) {
    const query: { [key: string]: string | string[] } = {};
    Object.keys(that.$route.query).map(
        (k) => (query[k] = that.$route.query[k])
    );
    query[key] = val;
    that.$router.replace({ query: query }).catch(() => { });
}
export async function GetCategories(id: string) {
    return (await axios.get<Category[]>('/api/relation/category/sub/' + id)).data ?? []
}

export async function UpsertCate(cate: Category) {
    return await (await axios.post<Category>("/api/blogs/category", cate)).data
}
