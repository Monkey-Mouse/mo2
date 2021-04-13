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
    settings?: UserSettings;
}
export interface DisplayBlogBrief extends BlogBrief {
    userLoad: boolean;
    userName: string;
    user?: User;
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
    perferDark?: string;
    themes?: string;
    bio?: string;
    github?: string;
    location?: string;
    github_id?: string;
    status?: string;
    home_img?: string;
    home_theme_dark?: "true" | 'false';
}
export interface User extends UserListData {
    email: string;
    roles?: string[];
    entityInfo?: EntityInfo;
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
    message?: string;
    iconClick?: (prop: InputProp) => void;
    onChange?: (c: any) => void;
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
export interface SubComment {
    authorProfile?: UserListData;
    aurhor?: string;
    content?: string;
    entity_info?: EntityInfo;
    id?: string;
    praise?: Praise;
    edit?: boolean;
}
export interface Praise {
    up?: number;
    down?: number;
    weight?: number;
}
export interface Comment extends SubComment {
    article?: string;
    subs?: SubComment[];
    tempC?: string;
    showSub?: boolean
}
export interface Count {
    count: number;
}
export interface Notification {
    operator_id: string;
    extra_message: string;
    create_time: string;
    processed: boolean;
}
export interface DisplayNotification extends Notification {
    user?: User;
}