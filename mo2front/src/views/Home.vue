<template>
  <div>
    <v-parallax
      src="https://cdn.vuetifyjs.com/images/parallax/material.jpg"
      height="200"
    >
      <v-row align="center" justify="center">
        <v-col class="text-center" cols="12">
          <h1 class="display-1 font-weight-thin mb-4">MO2</h1>
          <h4 class="subheading">Monkey ❤ Mouse</h4>
        </v-col>
      </v-row>
    </v-parallax>
    <v-container>
      <blog-time-line-list v-if="!firstloading" :blogs="blogs" />
      <blog-skeleton v-else :num="3" />
    </v-container>
  </div>
</template>

<script lang="ts">
import { BlogBrief } from "@/models";
import { GetArticles } from "@/utils";
import axios from "axios";
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
export default class Home extends Vue {
  blogs: BlogBrief[] = [];
  loading = true;
  firstloading = true;
  page = 0;
  pagesize = 5;
  nomore = false;
  created() {
    GetArticles({
      page: this.page++,
      pageSize: this.pagesize,
      draft: false,
    }).then((val) => {
      this.addMore(val);
      this.firstloading = false;
    });
  }
  addMore(val: BlogBrief[]) {
    if (val.length < this.pagesize) {
      this.nomore = true;
    }
    for (let index = 0; index < val.length; index++) {
      const element = val[index];
      this.blogs.push(element);
    }
    this.loading = false;
  }
  public ReachedButtom() {
    if (this.loading === false && !this.nomore) {
      this.loading = true;
      GetArticles({
        page: this.page++,
        pageSize: this.pagesize,
        draft: false,
      }).then((val) => {
        this.addMore(val);
      });
    }
  }
}
</script>