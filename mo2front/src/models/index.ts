export interface BlogBrief {
    id: string;
    title: string;
    cover?: string;
    rate?: number;
    description?: string;
    entityInfo: EntityInfo;
    authorId: string;
    score_sum?: number;
    score_num?: number;
    project_id?: string;
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
    y_doc?: string;
    is_y_doc?: boolean;
    y_token?: string;
}
export const BlankBlog: Blog = {
    title: "",
    content: "",
    id: "",
    cover: '',
    entityInfo: {},
    authorId: ''
}
export interface EntityInfo {
    createTime?: string;
    updateTime?: string;
    is_deleted?: boolean;
}
export interface BlogUpsert {
    id?: string;
    content?: string;
    cover?: string;
    description?: string;
    keyWords?: string[];
    title?: string;
    categories?: string[];
    y_doc?: string;
    is_y_doc?: boolean;
    y_token?: string;
    authorId?: string;
    project_id?:string;
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
type Options = Array<Option>|Array<string>;
export interface InputProp {
    errorMsg: { [name: string]: string };
    col?: number;
    type?: string;
    icon?: string;
    label?: string;
    default: any;
    accept?: string;
    options?: Options;
    showAvatar?: boolean;
    message?: string;
    multiple?:boolean;
    input?:string;
    iconClick?: (prop: InputProp) => void;
    onChange?: (c: any) => void;
    loading?:boolean;
    filter?:(item: any, queryText: string, itemText: string) => void
}
export interface Option {
    text: string;
    value: string;
    avatar?:string;
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
export interface Project {
    EntityInfo?: EntityInfo;
    ID?: string;
    Name?: string;
    Tags?: string[];
    OwnerID?: string;
    ManagerIDs?: string[];
    MemberIDs?: string[];
    BlogIDs?: string[];
    Description?: string;
    Avatar?: string;
}