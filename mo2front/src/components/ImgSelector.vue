<template>
  <v-container>
    <v-row class="grey--text">
      <v-col>{{ title }}</v-col>
      <v-spacer />
      <v-icon>{{ appendIcon }}</v-icon>
    </v-row>
    <v-row>
      <v-col
        v-for="(img, n) in imgs"
        :key="n"
        class="d-flex child-flex"
        cols="4"
      >
        <v-img
          @click="imgClick(img, n)"
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
  @Prop()
  title: string;
  @Prop()
  appendIcon: string;
  prevImg = 0;
  created() {
    for (let index = 0; index < this.imgs.length; index++) {
      const element = this.imgs[index];
      if (element.active) {
        this.prevImg = index;
      }
    }
  }
  imgClick(img: { src: string; active: boolean }, n: number) {
    if (this.multiple) {
      img.active = !img.active;
      return;
    }
    this.imgs[this.prevImg].active = false;
    img.active = true;
    this.prevImg = n;
    this.$emit("imgselect", img.src);
  }
}
</script>
<style>
.bordered {
  border: 2px solid #3298dc;
}
</style>
