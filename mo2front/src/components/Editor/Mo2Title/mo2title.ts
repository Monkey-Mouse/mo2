import { Node } from '@tiptap/core'

export const Mo2Title = Node.create({
    name: 'title',
    content: 'inline*',
    parseHTML() {
        return [{
            tag: 'h1',
        }];
    },
    renderHTML() {
        return ['h1', 0];
    },
    priority: 0
})