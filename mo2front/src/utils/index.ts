import Vue from '*.vue';
import { User, ApiError, InputProp, Option } from '../models/index'
import { AxiosError } from 'axios';
import * as qiniu from 'qiniu-js';
import { VuetifyThemeVariant } from 'vuetify/types/services/theme';
import { GetUploadToken, SearchUser, UpdateUserInfo } from './api';
import { LazyExecutor } from './lazy-executor';
export * from './api'
export * from './autoloader'
export * from './lazy-executor'
export function ElementInViewport(el: HTMLElement) {
    var top = el.offsetTop;
    var left = el.offsetLeft;
    var width = el.offsetWidth;
    var height = el.offsetHeight;

    while (el.offsetParent) {
        el = el.offsetParent as HTMLElement;
        top += el.offsetTop;
        left += el.offsetLeft;
    }

    return (
        top < (window.pageYOffset + window.innerHeight) &&
        left < (window.pageXOffset + window.innerWidth) &&
        (top + height) > window.pageYOffset &&
        (left + width) > window.pageXOffset
    );
}

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

export function GetInitials(name: string) {
    let rgx = new RegExp(/(\p{L}{1})\p{L}+/, 'gu');

    let initials = [...name.matchAll(rgx)] || [];

    return (
        (initials.shift()?.[1] || '') + (initials.pop()?.[1] || '')
    ).toUpperCase();
}
export function GetErrorMsg(apiError: any) {
    const err = (apiError as AxiosError<ApiError>);
    try {
        return err.response.data.reason
    } catch (error) {
        return 'Unknown Error'
    }
}
export function getRandomColor() {
    var letters = '456789ABCD';
    var color = '#';
    for (var i = 0; i < 6; i++) {
        color += letters[Math.floor(Math.random() * letters.length)];
    }
    return color;
}
export function makeid(length) {
    var result = [];
    var characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    var charactersLength = characters.length;
    for (var i = 0; i < length; i++) {
        result.push(characters.charAt(Math.floor(Math.random() *
            charactersLength)));
    }
    return result.join('');
}
export var globaldic: any = {};


export const AdminRole = "GeneralAdmin"
export const UserRole = "OrdinaryUser"
export const AnonymousRole = "Anonymous"
export function timeout(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}
export async function addQuery(that: Vue, key: string, val: string | string[]) {
    const query: { [key: string]: string | string[] } = {};
    Object.keys(that.$route.query).map(
        (k) => (query[k] = that.$route.query[k])
    );
    query[key] = val;
    that.$router.replace({ query: query }).catch(() => { });
}
interface App {
    refresh: boolean,
    showLogin: () => void,
    Prompt(msg: string,
        timeout: number): void,
    isUser: boolean,
    showGroup: boolean
}
var app: App;
export function SetApp(params: App) {
    app = params;
}
export function NewGroup() {
    if (!app.isUser) {
        Prompt("Please login first!", 5000)
        ShowLogin();
        return;
    }
    app.showGroup = true
}
export function ShowLogin() {
    app.showLogin()
}
export function Prompt(msg: string, timeout: number) {
    app.Prompt(msg, timeout)
}
export function ShowRefresh() {
    app.refresh = true;
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
            }).catch(err => reject(err))
        })
        promises.push(promise)
    }

    try {
        await Promise.all(promises)
    } catch (error) {
        Prompt(GetErrorMsg(error), 5000);
    }

}
export function LoginBeforeNav(to, from, next) {
    if (!app) {
        next()
        return
    }
    if (app.isUser) {
        next()
    } else {
        Prompt("Please login first!", 5000)
        ShowLogin()
    }
}
export function GetTheme() {
    return JSON.parse(
        localStorage.getItem("darkTheme")
    ) as boolean;
}
export function SetTheme(dark: boolean, that: Vue, themes?: { light: VuetifyThemeVariant, dark: VuetifyThemeVariant }, user?: User) {
    that.$vuetify.theme.dark = dark;
    localStorage.setItem("darkTheme", String(that.$vuetify.theme.dark));
    if (themes) {
        localStorage.setItem("themes", JSON.stringify(themes));
    }
    if (user && user.roles && user.roles.indexOf(UserRole) > -1) {
        if (!user.settings) {
            user.settings = {};
        }
        user.settings.perferDark = localStorage.getItem("darkTheme");
        user.settings.themes = localStorage.getItem("themes");
        UpdateUserInfo(user);
    }
}
export function SetThemeColors(that: Vue, themes?: { light: VuetifyThemeVariant, dark: VuetifyThemeVariant }) {
    for (const k in themes.dark) {
        that.$vuetify.theme.themes.dark[k] = themes.dark[k]
    }
    for (const k in themes.light) {
        that.$vuetify.theme.themes.light[k] = themes.light[k]
    }
}


export function ShareToQQ(param: { title: string, pic: string, summary: string, desc: string }) {
    window.open(`https://sns.qzone.qq.com/cgi-bin/qzshare/cgi_qzshare_onekey?url=${encodeURIComponent(document.location.toString())}&sharesource=qzone&title=${param.title}&pics=${param.pic}&summary=${param.summary}`, "_blank")
}
export function GithubOauth() {
    window.location.replace("https://github.com/login/oauth/authorize?client_id=c9cb620eaea6bff97e5d")
}
export function GenerateTOC() {
    var toc = "";
    var level = 0;
    var i = 0;
    let first = true;
    const processFunc = function (str, openLevel, attrs, titleText, closeLevel) {
        if (openLevel != closeLevel) {
            return str;
        }
        // openLevel -= 1;
        if (openLevel > level) {
            if (!first) {
                toc = toc.slice(0, toc.length - 9) + `
                <b
                    style="float: right;display:inline-block"
                    class="v-icon notranslate v-icon--link mdi mdi-chevron-left"
                    onclick="this.parentElement.parentElement.className==='active'?this.parentElement.parentElement.className='':this.parentElement.parentElement.className='active';event.preventDefault()"
                >
                </b>` + toc.slice(toc.length - 9)
            } else first = false;

            toc += (new Array(openLevel - level + 1)).join("<ul>");
        } else if (openLevel < level) {
            toc += (new Array(level - openLevel + 1)).join("</ul>");
        }

        level = parseInt(openLevel);

        var anchor = titleText.replace(/ /g, "_") + Date.now() + i++;
        toc += "<li><a href=\"#" + anchor + "\">" + `<div>${titleText}</div>`
            + "</a></li>";

        return "<h" + (openLevel) + attrs + ` id="${anchor}" anchor class="anchor h">`
            + titleText + "</h" + closeLevel + ">";
    }
    document.getElementById("titleContainer").innerHTML =
        document.getElementById("titleContainer").innerHTML.replace(
            /<h([\d])([^>]*)>([^<]+)<\/h([\d])>/gi,
            processFunc
        );
    document.getElementById("contents").innerHTML =
        document.getElementById("contents").innerHTML.replace(
            /<h([\d])([^>]*)>([^<]+)<\/h([\d])>/gi,
            processFunc
        );

    if (level) {
        toc += (new Array(level + 1)).join("</ul>");
    }
    const tocE = document.getElementById("toc");
    tocE.innerHTML += toc;
    let prevNode: Element = null;
    let prev: Element = null;
    setTimeout(() => {
        const hs = document.querySelectorAll("h1, h2, h3, h4, h5, h6")
        window.addEventListener('scroll', () => {
            for (const i of hs) {
                if (ElementInViewport(i as HTMLElement)) {
                    if (prev === i) {
                        return;
                    }
                    prev = i;
                    if (prevNode) {
                        prevNode.className = '';
                    }
                    const ae = tocE.querySelector('a[href="#' + i.id + '"]')
                    if (!ae) {
                        return;
                    }
                    const node = ae.parentElement
                    node.classList.add('active')
                    if (node.previousElementSibling && node.previousElementSibling.previousElementSibling) {
                        node.previousElementSibling.className = '';
                        node.previousElementSibling.previousElementSibling.className = '';
                    }
                    prevNode = node;
                    let n = node;
                    while (n.parentElement && n.parentElement.tagName === 'UL') {
                        n = n.parentElement;
                        if (n.previousElementSibling) {

                            n.previousElementSibling.className = 'active';
                        }
                    }
                    break;
                }
            }
        })
    }, 100);
};
export const BuildOnUserChange = (lazySearcher: LazyExecutor, dic: { [key: string]: Option }) => (a: { input: string; val: InputProp; cu: string[] }) => {
    a.val.input = a.input;
    if (!a.input) {
        return;
    }
    lazySearcher.Execute(() => {
        a.val.loading = true;
        SearchUser({ page: 0, pagesize: 5, query: a.input }).then((re) => {
            const set = new Set<Option>(a.cu.map((id) => dic[id]));
            re.map((u) => {
                const v = {
                    text: u.userName,
                    value: u.id,
                    avatar: u["settings.avatar"],
                    email: u.email,
                };
                dic[u.id] = v;
                set.add(v);
            });
            a.val.options = new Array(...set);
            a.val.loading = false;
        });
    });
}

export const UserFilter = (item: any, queryText: string, itemText: string) => {
    return (
        (item.email + item.text)
            .toLocaleLowerCase()
            .indexOf(queryText.toLocaleLowerCase()) > -1
    );
}