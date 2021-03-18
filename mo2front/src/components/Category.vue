<template>
  <v-container class="fill-height">
    <MO2Dialog
      :validator="validator"
      :inputProps="inputProps"
      :show.sync="addCate"
      title="添加集合"
      confirmText="确认"
      :confirm="confirm"
    />
    <MO2Dialog
      :validator="validator"
      :inputProps="inputPropsEdit"
      :show.sync="editCate"
      title="添加集合"
      confirmText="确认"
      :confirm="confirmEdit"
    />
    <v-row>
      <v-col>
        <v-breadcrumbs :items="items">
          <template v-slot:divider>
            <v-icon>mdi-forward</v-icon>
          </template>
          <template v-slot:item="{ item }">
            <v-breadcrumbs-item
              :disabled="item.id === items[items.length - 1].id"
            >
              <a @click="gotoLevel(item)">{{ item.text.toUpperCase() }}</a>
            </v-breadcrumbs-item>
          </template>
        </v-breadcrumbs>
      </v-col>

      <v-btn
        v-if="own"
        class="mt-5 mr-3"
        small
        @click="addCate = true"
        fab
        color="primary"
      >
        <v-icon> mdi-plus </v-icon>
      </v-btn>
    </v-row>
    <v-row v-if="loading" justify="center">
      <v-sheet v-for="i in 4" :key="i">
        <v-skeleton-loader
          class="ma-4"
          min-width="300"
          type="card@2"
        ></v-skeleton-loader>
      </v-sheet>
    </v-row>
    <v-row v-else-if="cate.length > 0">
      <v-col :key="i" v-for="(c, i) in cate">
        <v-hover v-slot="{ hover }">
          <v-card
            elevation="10"
            class="mx-auto"
            min-width="300"
            max-width="344"
          >
            <v-card-text>
              <div>mo2 category</div>
              <p class="display-1 text--primary">{{ c.name }}</p>
            </v-card-text>
            <v-card-actions>
              <v-chip class="ma-2" color="secondary" label text-color="white">
                <v-icon left> mdi-label </v-icon>
                Category
              </v-chip>
            </v-card-actions>

            <v-expand-transition>
              <v-card
                elevation="10"
                v-if="hover"
                class="transition-fast-in-fast-out v-card--reveal"
                style="height: 100%"
                color="primary"
                @click="nextCate(c)"
              >
                <v-card-text class="pb-0">
                  <p class="display-1 text--primary">Click to enter</p>
                  <!-- <p class="display-1 text--primary">Contains:</p> -->
                  <p class="display-1 text--secondary">
                    {{ c.name }}
                  </p>
                </v-card-text>
                <v-card-actions v-if="own">
                  <v-spacer />
                  <v-btn
                    @click="edit(c)"
                    v-on:click.prevent
                    v-on:click.stop
                    color="accent"
                    >Edit</v-btn
                  >
                  <v-btn
                    color="error"
                    v-on:click.prevent
                    v-on:click.stop
                    @click="deleteCate(c)"
                    >delete</v-btn
                  >
                </v-card-actions>
              </v-card>
            </v-expand-transition>
          </v-card>
        </v-hover>
      </v-col>
    </v-row>
    <nothing
      v-else-if="blogs.length === 0"
      :btnText="own ? 'Create New' : ''"
      @click="addCate = true"
    />
    <v-container>
      <v-divider></v-divider>
      <blog-time-line-list
        v-if="!loading"
        :blogs="blogs"
        :showNothing="false"
      />
      <blog-skeleton v-else :num="5" />
    </v-container>
  </v-container>
</template>
<script lang="ts">
import { Blog, Category, InputProp, User } from "@/models";
import {
  addQuery,
  DeleteCategories,
  GetCateBlogs,
  GetCategories,
  GetErrorMsg,
  UpsertCate,
} from "@/utils";
import Vue from "vue";
import Component from "vue-class-component";
import { Prop } from "vue-property-decorator";
import { required } from "vuelidate/lib/validators";
import MO2Dialog from "./MO2Dialog.vue";
import Nothing from "./NothingHere.vue";
import BlogTimeLineList from "../components/BlogTimeLineList.vue";
import BlogSkeleton from "../components/BlogTimeLineSkeleton.vue";
@Component({
  components: { MO2Dialog, Nothing, BlogTimeLineList, BlogSkeleton },
})
export default class Mo2Category extends Vue {
  @Prop()
  user!: User;
  @Prop()
  own: boolean;
  items: { text: string; id: string }[] = [];
  cate: Category[] = [];
  lev = 1;
  parentId = "";
  validator = {
    name: {
      required: required,
    },
  };
  inputProps: { [name: string]: InputProp } = {
    name: {
      errorMsg: {
        required: "集合名不可为空",
      },
      label: "Name",
      default: "",
      icon: "mdi-folder",
      col: 12,
      type: "text",
    },
  };
  inputPropsEdit: { [name: string]: InputProp } = {
    name: {
      errorMsg: {
        required: "集合名不可为空",
      },
      label: "Name",
      default: "",
      icon: "mdi-folder",
      col: 12,
      type: "text",
    },
  };
  addCate = false;
  editCate = false;
  loading = true;
  ec: Category = {};
  blogs: Blog[] = [];
  async confirm({ name }: { name: string }) {
    try {
      const data = await UpsertCate({
        name: name,
        parent_id: this.parentId === "" ? this.user.id : this.parentId,
      });
      this.cate.push(data);
      return { err: "", pass: true };
    } catch (error) {
      return { err: GetErrorMsg(error), pass: false };
    }
  }
  async edit(c: Category) {
    this.inputPropsEdit.name.default = c.name;
    this.editCate = true;
    this.ec = c;
  }
  async deleteCate(c: Category) {
    await DeleteCategories([c.id]);
    this.cate.splice(this.cate.indexOf(c));
  }
  async confirmEdit({ name }: { name: string }) {
    try {
      const c = this.ec;
      c.name = name;
      const data = await UpsertCate(c);
      c.name = data.name;
      return { err: "", pass: true };
    } catch (error) {
      return { err: GetErrorMsg(error), pass: false };
    }
  }
  gotoLevel(item: { text: string; id: string }) {
    this.loading = true;
    const pos = this.items.indexOf(item);
    this.lev = pos + 1;
    this.items = this.items.slice(0, pos + 1);
    const c = this.items[pos];
    this.parentId = c.id;
    addQuery(this, "cur", c.id);
    this.loadData(c.id);
  }
  created() {
    this.items.push({
      text: "Root",
      id: this.user.id,
    });
    const urlParams = new URLSearchParams(window.location.search);
    const cur = urlParams.get("cur");
    if (cur) {
      this.items = [];
      for (
        let index = 0;
        index < Object.keys(this.$route.query).length;
        index++
      ) {
        const element = this.$route.query[`lvid${index + 1}`] as string;
        const name = this.$route.query[`lv${index + 1}`] as string;
        this.items.push({ text: name, id: element });
        if (element === cur) {
          break;
        }
      }
      this.lev = this.items.length + 1;
      this.loadData(cur);
      return;
    }
    this.setCol("root", this.user.id);
    this.loadData(this.user.id);
  }
  loadData(id: string) {
    this.loading = true;
    GetCategories(id).then((data) => {
      this.cate = data;
      this.loading = false;
    });
    GetCateBlogs(id).then((data) => {
      this.blogs = data;
    });
  }
  setCol(name: string, id: string) {
    // addQuery(this, `lvid1`, this.user.id);
    addQuery(this, `lv${this.lev}`, name);
    addQuery(this, `lvid${this.lev}`, id);
    this.lev++;
    addQuery(this, "cur", id);
  }
  nextLev(c: Category) {
    this.parentId = c.id;
    this.loadData(c.id);
    const name = c.name;
    this.items.push({
      text: name,
      id: c.id,
    });
    this.setCol(name, c.id);
  }
  nextCate(c: Category) {
    this.loading = true;
    this.nextLev(c);
  }
}
</script>
<style>
.v-card--reveal {
  bottom: 0;
  opacity: 1 !important;
  position: absolute;
  width: 100%;
}
</style>