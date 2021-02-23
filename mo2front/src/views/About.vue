<template>
  <img-selector :imgs="imgs" />
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import { Prop } from "vue-property-decorator";
import ImgSelector from "../components/ImgSelector.vue";
@Component({
  components: { ImgSelector },
})
export default class About extends Vue {
  imgs: { src: string; active: boolean }[] = [
    { src: "https://picsum.photos/id/11/100/60", active: false },
    { src: "https://picsum.photos/id/11/100/60", active: false },
    { src: "https://picsum.photos/id/11/100/60", active: false },
    { src: "https://picsum.photos/id/11/100/60", active: false },
  ];
}
</script>
<style>
.bordered {
  border: 2px solid #3298dc;
}
</style>
