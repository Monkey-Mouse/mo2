export interface BlogBrief {
    id: string;
    title: string;
    cover: string;
    rate: number;
    description: string;
    entityInfo: EntityInfo;
    authorId: string;
}
export const BlankUser = {
    name: "",
    email: "",
    id: "",
    roles: [],
    settings: {
        avatar: ""
    },
    entityInfo: {}
}
export interface UserListData {
    id: string;
    name: string;
}
export interface DisplayBlogBrief extends BlogBrief {
    userLoad: boolean;
    userName: string;
}

export interface Blog extends BlogBrief {
    content: string;
}
export interface EntityInfo {
    createTime?: string;
    updateTime?: string;
}
export interface BlogUpsert {
    id?: string;
    content?: string;
    cover?: string;
    description?: string;
    keyWords?: string[];
    title?: string;
    categories?: string[]
}
export interface UserSettings {
    avatar?: string;
}
export interface User {
    id: string;
    name: string;
    email: string;
    roles?: string[];
    entityInfo?: EntityInfo;
    settings?: UserSettings;
}
export interface InputProp {
    errorMsg: { [name: string]: string };
    col?: number;
    type?: string;
    icon?: string;
    label?: string;
    default: any;
    accept?: string;
    options?: Option[];
    iconClick?: (prop: InputProp) => void;
}
export interface Option {
    text: string;
    value: string;
}
export interface ApiError {
    reason: string;
    time: string;
}
export interface ImgToken {
    token: string;
    file_key: string;
}

export interface Category {
    id?: string;
    parent_id?: string;
    name?: string;
}