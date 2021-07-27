<template>
  <div>
    <v-parallax
      src="https://cdn.mo2.leezeeyee.com/material.jpg~parallax"
      height="200"
    >
      <v-row align="center" justify="center">
        <v-col class="text-center" cols="12">
          <h1
            @click="$router.push('/about')"
            class="display-1 font-weight-thin mb-4 clickable"
          >
            MO2
          </h1>
          <h4 @click="$router.push('/about')" class="subheading clickable">
            Blog site build for everyone, built by everyone
          </h4>
        </v-col>
      </v-row>
    </v-parallax>
    <v-container>
      <blog-time-line-list v-if="!firstloading" :blogs="datalist" />
      <blog-skeleton v-if="loading" :num="pagesize" />
    </v-container>
  </div>
</template>

<script lang="ts">
import { BlogBrief } from "@/models";
import { GetArticles, AutoLoader, AddMore, ElmReachedBottom } from "../utils";
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
export default class Home extends Vue implements AutoLoader<BlogBrief> {
  datalist: BlogBrief[] = [];
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
    ElmReachedBottom(this, ({ page, pageSize }) =>
      GetArticles({
        page: page,
        pageSize: pageSize,
        draft: false,
      })
    );
  }
}
</script>