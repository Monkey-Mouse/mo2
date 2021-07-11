export function ParseQuery(query: { [key: string]: any }) {
    let queryStr = '?';
    const queryList: string[] = [];
    for (const key in query) {
        const element = query[key];
        if (element !== undefined && element !== null) {
            if (element instanceof Array) {
                element.forEach(e => {
                    queryList.push(`${key}=${e}`)
                });
            } else queryList.push(`${key}=${element}`)
        }
    }
    queryStr = queryStr + queryList.join('&');
    return queryStr
}