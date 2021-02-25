export interface BlogBrief {
    id: string;
    title: string;
    cover: string;
    rate: number;
    description: string;
    entityInfo: EntityInfo;
    authorId: string;
}
export interface Blog extends BlogBrief {
    content: string;
}
export interface EntityInfo {
    createTime: string;
    updateTime: string;
}
export interface BlogUpsert {
    id?: string;
    content?: string,
    cover?: string,
    description?: string,
    keyWords?: string[],
    title?: string
}
export interface User {
    id: string;
    name: string;
    email: string;
    description: string;
    createTime: string;
    site?: string;
    avatar: string;
    roles?: string[];
}
export interface InputProp {
    errorMsg: { [name: string]: string };
    col?: number;
    type?: string;
    icon?: string;
    label?: string;
    default: any;
    iconClick?: (prop: InputProp) => void;
}
export interface ApiError {
    reason: string;
    time: string;
}
export interface ImgToken {
    token: string;
    file_key: string;
}