<template>
  <v-container>
    <input
      type="file"
      ref="f"
      name="Upload Image"
      title="Upload Image"
      accept="image/*"
      style="display: none"
      @change="fileSelected"
    />
    <v-row v-show="connected || !this.$route.query['group']" justify="center">
      <v-col cols="12" lg="7" class="mo2editor">
        <div
          v-if="editor"
          v-show="this.editor.isFocused"
          class="menubar mt-4"
          style="z-index: 9999"
          :style="{ 'overflow-x': $vuetify.breakpoint.xsOnly ? 'auto' : '' }"
        >
          <div class="toolbar row">
            <span v-if="editor.isActive('table')">
              <button
                class="menubar__button"
                @click="editor.chain().focus().deleteTable().run()"
              >
                <v-icon>mdi-table-large-remove</v-icon>
              </button>
              <button
                class="menubar__button"
                @click="editor.chain().focus().addColumnBefore().run()"
              >
                <v-icon>mdi-table-column-plus-before</v-icon>
              </button>
              <button
                class="menubar__button"
                @click="editor.chain().focus().addColumnAfter().run()"
              >
                <v-icon>mdi-table-column-plus-after</v-icon>
              </button>
              <button
                class="menubar__button"
                @click="editor.chain().focus().deleteColumn().run()"
              >
                <v-icon>mdi-table-column-remove</v-icon>
              </button>
              <button
                class="menubar__button"
                @click="editor.chain().focus().addRowBefore().run()"
              >
                <v-icon>mdi-table-row-plus-before</v-icon>
              </button>
              <button
                class="menubar__button"
                @click="editor.chain().focus().addRowAfter().run()"
              >
                <v-icon>mdi-table-row-plus-after</v-icon>
              </button>
              <button
                class="menubar__button"
                @click="editor.chain().focus().deleteRow().run()"
              >
                <v-icon>mdi-table-row-remove</v-icon>
              </button>
              <button
                class="menubar__button"
                @click="editor.chain().focus().splitCell().run()"
              >
                <v-icon>mdi-table-split-cell</v-icon>
              </button>
              <button
                class="menubar__button"
                @click="editor.chain().focus().mergeCells().run()"
              >
                <v-icon>mdi-table-merge-cells</v-icon>
              </button>
            </span>
          </div>
        </div>
        <v-progress-circular
          v-if="isuploading"
          class="offset-10 offset-lg-7"
          style="position: fixed; z-index: 0"
          indeterminate
          color="amber"
        ></v-progress-circular>
        <v-switch
          hint="collab"
          persistent-hint
          v-else
          v-model="collab"
          class="offset-10 offset-lg-7"
          style="position: fixed; z-index: 0"
        ></v-switch>
        <bubble-menu
          v-if="editor"
          v-show="this.editor.isFocused"
          :editor="editor"
        >
          <v-btn-toggle
            dense
            :value="[]"
            :style="{ 'overflow-x': $vuetify.breakpoint.xsOnly ? 'auto' : '' }"
          >
            <v-btn
              @click="editor.chain().focus().toggleBold().run()"
              :class="{
                'v-item--active v-btn--active': editor.isActive('bold'),
              }"
            >
              <v-icon>mdi-format-bold</v-icon>
            </v-btn>
            <v-btn
              @click="editor.chain().focus().toggleItalic().run()"
              :class="{
                'v-item--active v-btn--active': editor.isActive('italic'),
              }"
            >
              <v-icon>mdi-format-italic</v-icon>
            </v-btn>
            <v-btn
              @click="editor.chain().focus().toggleCode().run()"
              :class="{
                'v-item--active v-btn--active': editor.isActive('code'),
              }"
            >
              <v-icon>mdi-code-tags</v-icon>
            </v-btn>
            <v-btn
              @click="editor.chain().focus().toggleStrike().run()"
              :class="{
                'v-item--active v-btn--active': editor.isActive('strike'),
              }"
            >
              <v-icon>mdi-format-strikethrough</v-icon>
            </v-btn>
            <v-btn
              @click="editor.chain().focus().toggleUnderline().run()"
              :class="{
                'v-item--active v-btn--active': editor.isActive('underline'),
              }"
            >
              <v-icon>mdi-format-underline</v-icon>
            </v-btn>
            <v-btn
              @click="editor.chain().focus().setTextAlign('left').run()"
              :class="{
                'v-item--active v-btn--active': editor.isActive({
                  textAlign: 'left',
                }),
              }"
            >
              <v-icon>mdi-format-align-left</v-icon>
            </v-btn>
            <v-btn
              @click="editor.chain().focus().setTextAlign('center').run()"
              :class="{
                'v-item--active v-btn--active': editor.isActive({
                  textAlign: 'center',
                }),
              }"
            >
              <v-icon>mdi-format-align-center</v-icon>
            </v-btn>
            <v-btn
              @click="editor.chain().focus().setTextAlign('right').run()"
              :class="{
                'v-item--active v-btn--active': editor.isActive({
                  textAlign: 'right',
                }),
              }"
            >
              <v-icon>mdi-format-align-right</v-icon>
            </v-btn>
            <v-btn
              @click="editor.chain().focus().setTextAlign('justify').run()"
              :class="{
                'v-item--active v-btn--active': editor.isActive({
                  textAlign: 'justify',
                }),
              }"
            >
              <v-icon>mdi-format-align-justify</v-icon>
            </v-btn>
            <form
              class="menububble__form"
              v-if="linkMenuIsActive"
              @submit.prevent="setLinkUrl(linkUrl)"
            >
              <input
                class="menububble__input"
                type="text"
                v-model="linkUrl"
                placeholder="https://"
                ref="linkInput"
                @keydown.esc="hideLinkMenu"
              />
              <button
                class="menububble__button"
                @click="setLinkUrl(null)"
                type="button"
              >
                <v-icon>mdi-close</v-icon>
              </button>
            </form>
            <v-btn
              v-else
              @click="showLinkMenu(editor.getMarkAttributes('link'))"
              :class="{
                'v-item--active v-btn--active': editor.isActive('link'),
              }"
            >
              <v-icon>mdi-link-variant-plus</v-icon>
            </v-btn>
          </v-btn-toggle>
        </bubble-menu>
        <floating-menu
          v-if="editor"
          v-show="this.editor.isFocused"
          :editor="editor"
        >
          <v-btn-toggle
            dense
            :value="[]"
            :style="{ 'overflow-x': $vuetify.breakpoint.xsOnly ? 'auto' : '' }"
          >
            <v-btn
              small
              @click="editor.chain().focus().toggleHeading({ level: 1 }).run()"
            >
              <v-icon>mdi-format-header-1</v-icon>
            </v-btn>
            <v-btn
              small
              @click="editor.chain().focus().toggleHeading({ level: 2 }).run()"
            >
              <v-icon>mdi-format-header-2</v-icon>
            </v-btn>
            <v-btn
              small
              @click="editor.chain().focus().toggleHeading({ level: 3 }).run()"
            >
              <v-icon>mdi-format-header-3</v-icon>
            </v-btn>
            <v-btn
              small
              @click="editor.chain().focus().toggleHeading({ level: 4 }).run()"
            >
              <v-icon>mdi-format-header-4</v-icon>
            </v-btn>
            <v-btn
              small
              @click="
                editor
                  .chain()
                  .focus()
                  .insertTable({ rows: 3, cols: 3, withHeaderRow: true })
                  .run()
              "
            >
              <v-icon>mdi-table-plus</v-icon>
            </v-btn>
            <v-btn
              small
              @click="editor.chain().focus().toggleBulletList().run()"
            >
              <v-icon>mdi-format-list-bulleted</v-icon>
            </v-btn>
            <v-btn
              small
              @click="editor.chain().focus().toggleOrderedList().run()"
            >
              <v-icon>mdi-format-list-numbered</v-icon>
            </v-btn>
            <v-btn small @click="editor.chain().focus().toggleTaskList().run()">
              <v-icon>mdi-format-list-checks</v-icon>
            </v-btn>
            <v-btn
              small
              @click="editor.chain().focus().toggleCodeBlock().run()"
            >
              <v-icon>mdi-code-json</v-icon>
            </v-btn>
            <v-btn
              small
              @click="editor.chain().focus().setHorizontalRule().run()"
            >
              <v-icon>mdi-arrow-split-horizontal</v-icon>
            </v-btn>
            <v-btn small @click="$refs.f.click()">
              <v-icon>mdi-image-plus</v-icon>
            </v-btn>
          </v-btn-toggle>
        </floating-menu>
        <div class="mo2content" spellcheck="false">
          <editor-content v-if="editor" :editor="editor" />
          <div v-if="editor" class="character-count grey--text mt-16">
            {{ editor.getCharacterCount() }} characters
            <span v-if="this.connected"> ,{{ this.userNum }} co-editors</span>
          </div>
        </div>
      </v-col>
    </v-row>
    <v-row v-if="!(connected || !this.$route.query['group'])" justify="center">
      <v-col align-self="center" class="text-center">
        <v-progress-circular
          size="128"
          indeterminate
          color="primary"
        ></v-progress-circular
      ></v-col>
    </v-row>
    <v-row
      v-if="!(connected || !this.$route.query['group'])"
      class="text-center"
      justify="center"
    >
      <v-col align-self="center"> Connecting to collaborators</v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
/* eslint-disable @typescript-eslint/camelcase */
import Vue from "vue";
import Component from "vue-class-component";
import {
  Editor,
  EditorContent,
  FloatingMenu,
  BubbleMenu,
  VueNodeViewRenderer,
} from "@tiptap/vue-2";
import { PasteHandlerExt } from "./PasteHandlerExt/pasteHandlerExt";
import Placeholder from "@tiptap/extension-placeholder";
import Paragraph from "@tiptap/extension-paragraph";
import Text from "@tiptap/extension-text";
import BlockQuote from "@tiptap/extension-blockquote";
import Italic from "@tiptap/extension-italic";
import Bold from "@tiptap/extension-bold";
import Strike from "@tiptap/extension-strike";
import Heading from "@tiptap/extension-heading";
import BulletList from "@tiptap/extension-bullet-list";
import OrderedList from "@tiptap/extension-ordered-list";
import ListItem from "@tiptap/extension-list-item";
import Doc from "@tiptap/extension-document";
import Code from "@tiptap/extension-code";
import History from "@tiptap/extension-history";
import Gapcursor from "@tiptap/extension-gapcursor";
import Table from "@tiptap/extension-table";
import TableRow from "@tiptap/extension-table-row";
import TableCell from "@tiptap/extension-table-cell";
import TableHeader from "@tiptap/extension-table-header";
import Typography from "@tiptap/extension-typography";
import CharacterCount from "@tiptap/extension-character-count";
import Image from "@tiptap/extension-image";
import Dropcursor from "@tiptap/extension-dropcursor";
import Link from "@tiptap/extension-link";
import TaskList from "@tiptap/extension-task-list";
import TaskItem from "@tiptap/extension-task-item";
import TextAlign from "@tiptap/extension-text-align";
import Hr from "@tiptap/extension-horizontal-rule";
import Underline from "@tiptap/extension-underline";
import Collaboration from "@tiptap/extension-collaboration";
import * as Y from "yjs";
import { WebrtcProvider } from "y-webrtc";
import { WebsocketProvider } from "y-websocket";
import { Prop, Watch } from "vue-property-decorator";
import { getRandomColor, Prompt, SetBlogType, timeout } from "@/utils";
import CodeBlockLowlight from "@tiptap/extension-code-block-lowlight";
import CodeBlockComponent from "./Lowlight/CodeBlockComponent.vue";
import CollaborationCursor from "@tiptap/extension-collaboration-cursor";

// load all highlight.js languages
import lowlight from "lowlight";
import { User } from "@/models";

@Component({
  components: {
    EditorContent,
    FloatingMenu,
    BubbleMenu,
  },
})
export default class MO2Editor extends Vue {
  @Prop()
  content?: string;
  @Prop()
  uploadImgs: (
    blobs: File[],
    callback: (imgprop: { src: string; alt?: string; title?: string }) => void
  ) => Promise<void>;
  @Prop()
  user: User;
  @Prop()
  ystate: string;
  // @Prop({default:false})
  // isydoc:boolean;
  isuploading = false;
  linkMenuIsActive = false;
  update = false;
  editable = true;
  load = false;
  editor: Editor = null;
  provider: WebsocketProvider = null;
  rtcProvider: WebrtcProvider = null;
  ydoc: Y.Doc = null;
  userNum = 0;
  connected = false;
  get collab() {
    return this.$route.query["group"] !== undefined;
  }
  set collab(v: boolean) {
    const c = this.GetHTML();
    if (v) {
      // eslint-disable-next-line @typescript-eslint/camelcase
      SetBlogType({
        y_doc: "",
        is_y_doc: true,
        id: this.$route.params["id"],
      })
        .then((d) => {
          this.initEditor(c, d.token, true);
          this.$router.replace(this.$route.fullPath + "?group=" + d.token);
          navigator.clipboard.writeText(window.location.href).then(() => {
            Prompt("合作编辑加入链接已复制到剪贴板！分享给别人即可", 10000);
          });
        })
        .catch(() => {
          this.$router.replace(this.$route.path);
          Prompt("初始化合作编辑失败！", 10000);
        });
    } else {
      // eslint-disable-next-line @typescript-eslint/camelcase
      SetBlogType({
        y_doc: "",
        is_y_doc: false,
        id: this.$route.params["id"],
      }).then(() => {
        this.$router.replace(this.$route.path).then(() => {
          this.initEditor(c);
        });
        Prompt("退出合作编辑！", 10000);
      });
    }
  }

  initEditor(content: string, group?: string, dontUseYstate?: boolean) {
    const exts = [
      Paragraph,
      Text,
      Doc,
      Italic,
      Bold,
      BulletList,
      OrderedList,
      Strike,
      BlockQuote,
      ListItem,
      Code,
      Gapcursor,
      // History,
      Typography,
      CharacterCount,
      Table.configure({
        resizable: true,
      }),
      TableRow,
      TableHeader,
      TableCell,
      Image,
      Dropcursor,
      Link,
      TaskList,
      TaskItem,
      TextAlign,
      Hr,
      Underline,
      PasteHandlerExt.configure({
        uploadImgs: this.uploadImages,
      }),
      Placeholder.configure({ showOnlyCurrent: false }),
      Heading.configure({ levels: [1, 2, 3, 4] }),
      CodeBlockLowlight.extend({
        addNodeView() {
          return VueNodeViewRenderer(CodeBlockComponent as any);
        },
      }).configure({ lowlight }),
    ];
    if (!content || content.length === 0) {
      content = "<h1></h1><p></p>";
    }
    if (this.editor) {
      this.dispose();
    }
    group = group ?? (this.$route.query["group"] as string);
    if (group) {
      if (!dontUseYstate) {
        content = "";
      }

      this.ydoc = new Y.Doc();
      if (this.ystate && !dontUseYstate) {
        console.log(this.ystate);
        Y.applyUpdateV2(
          this.ydoc,
          Uint8Array.from(atob(this.ystate), (c) => c.charCodeAt(0))
        );
      }
      this.rtcProvider = new WebrtcProvider(group, this.ydoc);
      this.provider = new WebsocketProvider(
        "wss://demos.yjs.dev",
        group,
        this.ydoc
      );
      this.provider.on("status", (event) => {
        this.connected = event.status === "connected";
      });
      exts.push(
        Collaboration.configure({
          document: this.ydoc,
        }),
        CollaborationCursor.configure({
          provider: this.provider,
          user: {
            name: this.user.name,
            color: getRandomColor(),
          },
          onUpdate: (users) => {
            this.userNum = users.length;
            return null;
          },
        })
      );
    } else exts.push(History);
    // eslint-disable-next-line @typescript-eslint/no-this-alias
    const that = this;
    this.editor = new Editor({
      content: content,
      extensions: exts,
      onUpdate() {
        that.update = true;
      },
    });
  }
  fileSelected() {
    this.uploadImages([...(this.$refs.f as HTMLInputElement).files]);
    (this.$refs.f as HTMLInputElement).value = "";
  }
  uploadImages(files: File[]) {
    this.isuploading = true;
    this.uploadImgs(files, this.editor.commands.setImage)
      .then(() => (this.isuploading = false))
      .catch(() => {
        this.isuploading = false;
      });
  }
  mounted() {
    this.initEditor(this.content);
    this.$emit("loaded", this);
    this.startAutoSave();
  }
  async startAutoSave() {
    console.log("[editor]: start auto save");
    // eslint-disable-next-line no-constant-condition
    while (true) {
      if (this.update) {
        this.$emit("autosave");
        this.update = false;
      }
      await timeout(5000);
    }
  }
  @Watch("content")
  contentSet() {
    if (this.content !== null && this.content !== undefined) {
      this.initEditor(this.content);
    }
  }
  dispose() {
    this.editor.destroy();
    this.ydoc?.destroy();
    this.rtcProvider?.disconnect();
    this.provider?.disconnect();
    this.rtcProvider?.destroy();
    this.provider?.destroy();
  }

  beforeDestroy() {
    this.dispose();
  }
  linkUrl = null;
  showLinkMenu(attrs) {
    this.linkUrl = attrs.href;
    this.linkMenuIsActive = true;
    this.$nextTick(() => {
      (this.$refs.linkInput as HTMLElement).focus();
    });
  }
  hideLinkMenu() {
    this.linkUrl = null;
    this.linkMenuIsActive = false;
  }
  setLinkUrl(url: string) {
    if (!url || !url.startsWith("http")) {
      this.editor.commands.unsetLink();
      this.hideLinkMenu();
      return;
    }
    this.editor.commands.setLink({ href: url });
    this.hideLinkMenu();
  }
  @Watch("editable")
  changeEditable() {
    this.editor.setOptions({
      editable: this.editable,
    });
  }
  public GetHTML() {
    return this.editor.getHTML() as string;
  }
  Uint8ToString(u8a) {
    let CHUNK_SZ = 0x8000;
    let c = [];
    for (let i = 0; i < u8a.length; i += CHUNK_SZ) {
      c.push(String.fromCharCode.apply(null, u8a.subarray(i, i + CHUNK_SZ)));
    }
    return c.join("");
  }
  public GetYDoc() {
    return this.ydoc
      ? btoa(this.Uint8ToString(Y.encodeStateAsUpdateV2(this.ydoc)))
      : "";
  }
}
</script>

<style lang="scss">
/* Give a remote user a caret */
.collaboration-cursor__caret {
  position: relative;
  margin-left: -1px;
  margin-right: -1px;
  border-left: 1px solid #0d0d0d;
  border-right: 1px solid #0d0d0d;
  word-break: normal;
  pointer-events: none;
}

/* Render the username above the caret */
.collaboration-cursor__label {
  position: absolute;
  top: -1.4em;
  left: -1px;
  font-size: 12px;
  font-style: normal;
  font-weight: 600;
  line-height: normal;
  user-select: none;
  color: #0d0d0d;
  padding: 0.1rem 0.3rem;
  border-radius: 3px 3px 3px 0;
  white-space: nowrap;
}
</style>