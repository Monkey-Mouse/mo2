import Vue from '*.vue';
import { User, ApiError, ImgToken, BlogBrief, BlogUpsert, Blog, UserListData, Category, Comment, SubComment, Count, Notification } from '@/models/index'
import axios, { AxiosError } from 'axios';
import * as qiniu from 'qiniu-js';
import { VuetifyThemeVariant } from 'vuetify/types/services/theme';
import router from '../router'
import { UpdateUserInfo } from './api';
export * from './api'
export * from './autoloader'
export * from './lazy-executor'

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
        if (err.response.status === 404) {
            router.push('/404')
        }
        return err.response.data.reason
    } catch (error) {
        return 'Unknown Error'
    }
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

var app: { refresh: boolean, showLogin: () => void };
export function SetApp(params: { refresh: boolean, showLogin: () => void }) {
    app = params;
}
export function ShowLogin() {
    app.showLogin()
}
export function ShowRefresh() {
    app.refresh = true;
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

    document.getElementById("contents").innerHTML =
        document.getElementById("contents").innerHTML.replace(
            /<h([\d])>([^<]+)<\/h([\d])>/gi,
            function (str, openLevel, titleText, closeLevel) {
                if (openLevel != closeLevel) {
                    return str;
                }
                openLevel -= 1;
                if (openLevel > level) {
                    toc += (new Array(openLevel - level + 1)).join("<ul>");
                } else if (openLevel < level) {
                    toc += (new Array(level - openLevel + 1)).join("</ul>");
                }

                level = parseInt(openLevel);

                var anchor = titleText.replace(/ /g, "_") + Date.now() + i++;
                toc += "<li><a href=\"#" + anchor + "\">" + titleText
                    + "</a></li>";

                return "<h" + (openLevel + 1) + ` id="${anchor}" class="anchor">`
                    + titleText + "</h" + closeLevel + ">";
            }
        );

    if (level) {
        toc += (new Array(level + 1)).join("</ul>");
    }

    document.getElementById("toc").innerHTML += toc;
};