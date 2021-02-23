<template>
  <v-container>
    <v-row>
      <v-col
        v-for="(img, n) in imgs"
        :key="n"
        class="d-flex child-flex"
        cols="4"
      >
        <v-img
          @click="imgClick(img)"
          :src="img.src"
          aspect-ratio="1"
          class="black lighten-2 clickable is-clickable"
          :class="img.active ? 'elevation-24 bordered' : ''"
        >
          <template v-slot:placeholder>
            <v-row class="fill-height ma-0" align="center" justify="center">
              <v-progress-circular
                indeterminate
                color="grey lighten-5"
              ></v-progress-circular>
            </v-row>
          </template>
        </v-img>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import { Prop } from "vue-property-decorator";
@Component({
  components: {},
})
export default class ImgSelector extends Vue {
  @Prop()
  multiple!: boolean;
  @Prop()
  imgs: { src: string; active: boolean }[];
  prevImg = { src: "", active: false };
  imgClick(img: { src: string; active: boolean }) {
    if (this.multiple) {
      if (img.active) {
        img.active = false;
      }
      return;
    }
    this.prevImg.active = false;
    img.active = true;
    this.prevImg = img;
    this.$emit("select", img.src);
  }
}
</script>
<style>
.bordered {
  border: 2px solid #3298dc;
}
</style>
