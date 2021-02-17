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