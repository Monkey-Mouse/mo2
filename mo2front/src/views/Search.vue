<template>
  <div>
    <v-container>
      <v-row>
        <v-col cols="12" lg="10" offset-lg="1">
          <h1
            @keydown="keyDown"
            ref="s"
            contenteditable="true"
            class="display-1 text-h3 pt-10 font-weight-thin ml-4 mr-4 mb-0 borderless"
          >
            {{ search }}
          </h1>
        </v-col>
      </v-row>
      <v-row
        ><v-col cols="12" lg="10" offset-lg="1"><v-divider /></v-col
      ></v-row>
      <blog-time-line-list v-if="!firstloading" :blogs="blogs" />
      <blog-skeleton v-if="loading" :num="pagesize" />
    </v-container>
  </div>
</template>

<script lang="ts">
import { BlogBrief } from "@/models";
import {
  GetArticles,
  BlogAutoLoader,
  AddMore,
  ElmReachedButtom,
} from "../utils";
import Vue from "vue";
import Component from "vue-class-component";
import BlogTimeLineList from "../components/BlogTimeLineList.vue";
import BlogSkeleton from "../components/BlogTimeLineSkeleton.vue";
@Component({
  components: {
    BlogTimeLineList,
    BlogSkeleton,
  },
})
export default class Search extends Vue implements BlogAutoLoader {
  blogs: BlogBrief[] = [];
  loading = true;
  firstloading = true;
  search = "";
  page = 0;
  pagesize = 5;
  nomore = false;
  i = 0;
  keyDown(event: KeyboardEvent) {
    // Number 13 is the "Enter" key on the keyboard
    if (event.key === "Enter") {
      // Cancel the default action, if needed
      event.preventDefault();

      // Trigger the button element with a click
      this.$router
        .replace("search?q=" + (this.$refs.s as HTMLElement).textContent)
        .catch();
      if (!this.loading) {
        this.init();
      }
    } else {
      this.i++;
      let j = this.i;

      setTimeout(() => {
        if (j === this.i && !this.loading) {
          this.$router
            .replace("search?q=" + (this.$refs.s as HTMLElement).textContent)
            .catch();
          this.init();
        }
      }, 500);
    }
  }
  created() {
    this.search = this.$route.query["q"] as string;
    this.init();
  }
  init() {
    this.page = 0;
    this.blogs = [];
    this.firstloading = true;
    this.loading = true;
    this.nomore = false;
    GetArticles({
      page: this.page++,
      pageSize: this.pagesize,
      draft: false,
      search: this.$route.query["q"] as string,
    }).then((val) => {
      this.addMore(val);
      this.firstloading = false;
    });
  }
  addMore(val: BlogBrief[]) {
    AddMore(this, val);
  }
  mounted() {
    (this.$refs.s as HTMLElement).focus();
  }
  public ReachedButtom() {
    ElmReachedButtom(this, ({ page, pageSize }) =>
      GetArticles({
        page: page,
        pageSize: pageSize,
        draft: false,
        search: this.$route.query["q"] as string,
      })
    );
  }
}
</script>
<style scoped>
.borderless {
  outline: 0px solid transparent;
}
</style>