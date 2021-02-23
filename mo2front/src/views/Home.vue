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
import Vue from "vue";
import Component from "vue-class-component";
import BlogTimeLineList from "../components/BlogTimeLineList.vue";
@Component({
  components: {
    BlogTimeLineList,
  },
})
export default class Home extends Vue {
  blogs: BlogBrief[] = Array<BlogBrief>(10).fill({
    id: "string",
    title: "MO2",
    cover: "https://picsum.photos/500/300?image=40",
    rate: 4.3,
    description:
      "Lorem ipsum dolor sit amet, no nam oblique veritus. Commune scaevola imperdiet nec ut, sed euismod convenire principes at. Est et nobis iisque percipit, an vim zril disputando voluptatibus, vix an salutandi sententiae.",
    entityInfo: { createTime: "2021/2/9", updateTime: "2021/2/9" },
    author: "Leezeeyee",
  });
}
</script>