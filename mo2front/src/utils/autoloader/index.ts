import { AxiosError } from "axios";
import { reachedBottom } from './reachedBottom'

export function ReachedBottom(): boolean {
    return reachedBottom();
}
export interface AutoLoader<T> {
    datalist: T[];
    loading: boolean;
    firstloading: boolean;
    page: number;
    pagesize: number;
    nomore: boolean;
    ReachedButtom: () => void;
}
export function InitLoader<T>(loader: AutoLoader<T>) {
    loader.datalist = [];
    loader.page = 0;
    loader.firstloading = true;
    loader.nomore = false;
    loader.loading = true;
}
export function ElmReachedBottom<T>(elm: AutoLoader<T>, getMore: (query: { page: number; pageSize: number }) => Promise<T[]>) {
    if (elm.loading === false && !elm.nomore) {
        elm.loading = true;
        getMore({
            page: elm.page++,
            pageSize: elm.pagesize,
        }).then((val) => {
            try {
                AddMore(elm, val);
            } finally {
                elm.loading = false;
            }
        }).catch((err: AxiosError) => { elm.loading = false; });
    }
}
export function AddMore<T>(elm: AutoLoader<T>, val: T[]) {
    if (!val || val.length < elm.pagesize) {
        elm.nomore = true;
    }
    if (!val) {
        elm.loading = false;
        return
    }
    for (let index = 0; index < val.length; index++) {
        const element = val[index];
        elm.datalist.push(element);
    }
    elm.loading = false;
}