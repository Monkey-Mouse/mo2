<template>
  <blog-skeleton :num="10" />
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import { Prop } from "vue-property-decorator";
import BlogSkeleton from "../components/BlogTimeLineSkeleton.vue";
@Component({
  components: { BlogSkeleton },
})
export default class About extends Vue {}
</script>
<style>
.bordered {
  border: 2px solid #3298dc;
}
</style>
