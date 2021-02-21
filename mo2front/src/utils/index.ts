import { User, ApiError, ImgToken } from '@/models/index'
import axios, { AxiosError } from 'axios';
import * as qiniu from 'qiniu-js';

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
    throw new Error("Not implement yet");
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
export async function RegisterAsync(userInfo: { email: string, password: string, userName: string }) {
    return (await axios.post<User>('/api/accounts', userInfo)).data;
}
export async function LoginAsync(userInfo: { userNameOrEmail: string, password: string }) {
    return (await axios.post<User>('/api/accounts/login', userInfo)).data;
}
export function GetErrorMsg(apiError: any) {
    try {
        return (apiError as AxiosError<ApiError>).response.data.reason
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
                var ob = qiniu.upload(element, val.file_key, val.token);
                ob.subscribe(null, (err) => {
                    reject(err)
                }, res => {
                    callback({ src: '//qotwmtnjo.hn-bkt.clouddn.com/' + res.key })
                    resolve();
                })
            })
        })
        promises.push(promise)
    }
    await Promise.all(promises)

}