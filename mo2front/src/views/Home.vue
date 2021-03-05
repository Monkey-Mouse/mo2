<template>
  <div>
    <v-parallax
      src="https://cdn.vuetifyjs.com/images/parallax/material.jpg"
      height="200"
    >
      <v-row align="center" justify="center">
        <v-col class="text-center" cols="12">
          <h1 class="display-1 font-weight-thin mb-4">MO2</h1>
          <h4 class="subheading">
            Blog site build for everyone, built by everyone
          </h4>
        </v-col>
      </v-row>
    </v-parallax>
    <v-container>
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
} from "@/utils";
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
export default class Home extends Vue implements BlogAutoLoader {
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
    AddMore(this, val);
  }
  public ReachedButtom() {
    ElmReachedButtom(this, ({ page, pageSize }) =>
      GetArticles({
        page: page,
        pageSize: pageSize,
        draft: false,
      })
    );
  }
}
</script>