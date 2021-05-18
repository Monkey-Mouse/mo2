//#region code from @democrazy, limfx. site:www.limfx.pro


/**
 * 滚动条在Y轴上的滚动距离
 */
function getScrollTop(): number {
    let scrollTop = 0, bodyScrollTop = 0, documentScrollTop = 0;
    if (document.body) {
        bodyScrollTop = document.body.scrollTop;
    }
    if (document.documentElement) {
        documentScrollTop = document.documentElement.scrollTop;
    }
    scrollTop = (bodyScrollTop - documentScrollTop > 0) ? bodyScrollTop : documentScrollTop;
    return scrollTop as number;
}
/**
 * 文档的总高度
 */
function getScrollHeight(): number {
    var scrollHeight = 0, bodyScrollHeight = 0, documentScrollHeight = 0;
    let bSH = 0;
    if (document.body) {
        bSH = document.body.scrollHeight;
    }
    let dSH = 0;
    if (document.documentElement) {
        dSH = document.documentElement.scrollHeight;
    }
    scrollHeight = (bSH - dSH > 0) ? bSH : dSH;
    return scrollHeight;
}
/**
 * 浏览器视口的高度
 */
function getWindowHeight(): number {
    var windowHeight = 0;
    if (document.compatMode == "CSS1Compat") {
        windowHeight = document.documentElement.clientHeight;
    } else {
        windowHeight = document.body.clientHeight;
    }
    return windowHeight;
}
function checkVisible(elm: HTMLElement) {
    var rect = elm.getBoundingClientRect();
    //获取当前浏览器的视口高度，不包括工具栏和滚动条
    //document.documentElement.clientHeight兼容 Internet Explorer 8、7、6、5
    var viewHeight = Math.max(document.documentElement.clientHeight, window.innerHeight);
    //bottom top是相对于视口的左上角位置
    //bottom大于0或者top-视口高度小于0可见
    return !(rect.bottom < 0 || rect.top - viewHeight >= 0);
}
/**
 * 滚动到最底部
 */
export function reachedBottom(): boolean {
    let footer = document.getElementById('footer');
    if (footer) {
        return checkVisible(footer as HTMLElement);
    } else {
        if (Math.abs(getScrollTop() + getWindowHeight() - getScrollHeight()) <= 10) {
            return true;
        }
        return false;
    }
}
//#endregion