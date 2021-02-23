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
      <blog-time-line-list :blogs="blogs" />
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
@Component({
  components: {
    BlogTimeLineList,
  },
})
export default class Home extends Vue {
  blogs: BlogBrief[] = [];
  created() {
    GetArticles({ page: 0, pageSize: 10, draft: false }).then((val) => {
      this.blogs = val;
    });
  }
}
</script>