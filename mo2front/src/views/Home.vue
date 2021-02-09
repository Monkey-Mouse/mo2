<template>
  <blog-time-line-list :blogs="blogs" />
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
  blogs: BlogBrief[] = Array<BlogBrief>(100).fill({
    id: "string",
    title: "MO2",
    cover: "https://picsum.photos/500/300?image=40",
    rate: 4.3,
    description:
      "Lorem ipsum dolor sit amet, no nam oblique veritus. Commune scaevola imperdiet nec ut, sed euismod convenire principes at. Est et nobis iisque percipit, an vim zril disputando voluptatibus, vix an salutandi sententiae.",
    createTime: "2021/2/9",
    author: "Leezeeyee",
  });
  created() {}
}
</script>