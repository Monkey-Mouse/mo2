<template>
  <v-container>
    <input
      type="file"
      ref="f"
      accept="image/*"
      style="display: none"
      @change="fileSelected"
    />
    <v-row justify="center">
      <v-col cols="12" lg="7" class="mo2editor">
        <div v-if="editor" class="menubar mt-4" style="z-index: 9999">
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
          style="position: fixed; z-index: 999"
          indeterminate
          color="amber"
        ></v-progress-circular>
        <bubble-menu v-if="editor" :editor="editor">
          <v-btn-toggle dense :value="[]">
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
          </v-btn-toggle>
        </bubble-menu>
        <floating-menu v-if="editor" :editor="editor">
          <v-btn-toggle dense :value="[]">
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
            <v-btn
              small
              @click="editor.chain().focus().toggleCodeBlock().run()"
            >
              <v-icon>mdi-code-json</v-icon>
            </v-btn>
          </v-btn-toggle>
        </floating-menu>
        <div class="mo2content" spellcheck="false">
          <editor-content v-if="editor" :editor="editor" />
        </div>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import {
  Editor,
  EditorContent,
  FloatingMenu,
  BubbleMenu,
  VueNodeViewRenderer,
} from "@tiptap/vue-2";
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
import { Prop, Watch } from "vue-property-decorator";
import { timeout } from "@/utils";
import CodeBlockLowlight from "@tiptap/extension-code-block-lowlight";
import CodeBlockComponent from "./Lowlight/CodeBlockComponent.vue";
// load all highlight.js languages
import lowlight from "lowlight";

let that: MO2Editor | any = {};
@Component({
  components: {
    EditorContent,
    FloatingMenu,
    BubbleMenu,
    // EditorMenuBar,
  },
})
export default class MO2Editor extends Vue {
  @Prop()
  content?: string;
  @Prop()
  uploadImgs: (
    blobs: File[],
    callback: (imgprop: { src: string }) => void
  ) => Promise<void>;
  isuploading = false;
  linkMenuIsActive = false;
  update = false;
  editable = true;
  load = false;
  editor: Editor = null;
  initEditor(content: string) {
    console.log("init");
    if (this.editor) {
      this.editor.commands.setContent(content ?? "<h1></h1>");
      return;
    }
    if (!content || content.length === 0) {
      content = "<h1></h1><p></p>";
    }
    this.editor = new Editor({
      extensions: [
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
        History,
        Typography,
        CharacterCount,
        Table.configure({
          resizable: true,
        }),
        TableRow,
        TableHeader,
        TableCell,
        Placeholder.configure({ showOnlyCurrent: false }),
        Heading.configure({ levels: [1, 2, 3, 4] }),
        CodeBlockLowlight.extend({
          addNodeView() {
            return VueNodeViewRenderer(CodeBlockComponent as any);
          },
        }).configure({ lowlight }),
      ],
      content: content ?? "<h1></h1>",
      onUpdate() {
        (that as MO2Editor).update = true;
      },
    });
  }
  fileSelected() {
    this.isuploading = true;
    // this.uploadImgs(
    //   [...(this.$refs.f as HTMLInputElement).files],
    //   this.editor.commands.image
    // )
    //   .then(() => (that.isuploading = false))
    //   .catch(() => {
    //     that.isuploading = false;
    //   });
    (this.$refs.f as HTMLInputElement).value = "";
  }
  mounted() {
    that = this;
    this.initEditor(this.content);
    this.$emit("loaded", this);
    this.startAutoSave();
  }
  async startAutoSave() {
    while (true) {
      if (this.update) {
        (that as MO2Editor).$emit("autosave");
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

  beforeDestroy() {
    this.editor.destroy();
  }
  linkUrl = null;
  showLinkMenu(attrs) {
    this.linkUrl = attrs.href;
    this.linkMenuIsActive = true;
    this.$nextTick(() => {
      (this.$refs.linkInput as any).focus();
    });
  }
  hideLinkMenu() {
    this.linkUrl = null;
    this.linkMenuIsActive = false;
  }
  setLinkUrl(command, url) {
    command({ href: url });
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
}
</script>