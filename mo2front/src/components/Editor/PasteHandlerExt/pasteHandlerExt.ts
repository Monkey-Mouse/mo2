import { Extension } from '@tiptap/core'
import { Plugin } from 'prosemirror-state'

export interface PasteExtOptions {
    uploadImgs?: (
        blobs: File[],
    ) => void;
}

export const PasteHandlerExt = Extension.create<PasteExtOptions>({
    defaultOptions: {
        uploadImgs: null
    },
    // Your code goes here.
    addProseMirrorPlugins() {
        return [
            new Plugin({
                props: {
                    handlePaste: (v, event, slice) => {
                        let items = (event.clipboardData)
                            .items;
                        let files = [];
                        for (let index = 0; index < items.length; index++) {
                            let item = items[index] as DataTransferItem;
                            if (
                                item.type === "text/html" &&
                                items.length >= index + 1 &&
                                items[index + 1].kind === "file"
                            ) {
                                index++;
                                continue;
                            }
                            if (item.kind === "file") {
                                let blob = item.getAsFile();
                                files = files.concat(blob);
                                // var reader = new FileReader();
                                // reader.onload = function(event){
                                //   console.log(event.target.result)}; // data url!
                                // reader.readAsDataURL(blob);
                            }
                        }
                        if (files.length === 0) {
                            return false;
                        }
                        (event as Event).stopImmediatePropagation();
                        (event as Event).stopPropagation();
                        (event as Event).preventDefault();
                        if (this.options.uploadImgs) this.options.uploadImgs(files)

                    },
                }
            })
        ]
    }
})