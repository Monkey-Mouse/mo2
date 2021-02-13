<template>
  <v-timeline :dense="this.$vuetify.breakpoint.smAndDown">
    <v-timeline-item
      v-for="(blog, i) in blogs"
      :key="i"
      :color="displayColors[i % displayColors.length]"
      large
    >
      <v-lazy
        :value="false"
        :options="{
          threshold: 0.5,
        }"
        min-height="200"
        transition="fade-transition"
      >
        <v-card class="mx-auto elevation-20" style="max-width: 400px">
          <v-row justify="space-between">
            <v-col sm="6" cols="12">
              <v-card-title>
                <div>
                  <div
                    :class="`${displayColors[i % displayColors.length]}--text`"
                    class="headline"
                  >
                    {{ blog.title }}
                  </div>
                  <div class="subtitle-1">{{ blog.author }}</div>
                  <div class="subtitle-2">{{ blog.createTime }}</div>
                </div>
              </v-card-title>
            </v-col>
            <v-col sm="6" cols="0">
              <v-img
                class="shrink ma-3"
                contain
                height="125px"
                :src="blog.cover"
                style="flex-basis: 125px"
              ></v-img>
            </v-col>
          </v-row>
          <v-divider dark></v-divider>
          <v-card-actions class="pa-4">
            Rate this
            <v-spacer></v-spacer>
            <span class="text--lighten-2 caption mr-2">
              ({{ blog.rate }})
            </span>
            <v-rating
              v-model="blog.rate"
              background-color="white"
              color="yellow accent-4"
              @change="rateChange(blog)"
              dense
              half-increments
              hover
              size="18"
            ></v-rating>
          </v-card-actions>
        </v-card>
      </v-lazy>
      <template v-slot:opposite>
        <v-lazy
          :value="false"
          :options="{
            threshold: 0.5,
          }"
          min-height="200"
          transition="fade-transition"
        >
          <div class="py-4">
            <div>
              {{ blog.description }}
            </div>
          </div>
        </v-lazy>
      </template>
    </v-timeline-item>
  </v-timeline>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import { Prop } from "vue-property-decorator";
import { Colors } from "vuetify/es5/util/colors";
import { colors } from "vuetify/lib";
import { BlogBrief } from "../models/index";
import { randomProperty, Copy } from "../utils/index";
@Component
export default class BlogTimeLineList extends Vue {
  @Prop()
  blogs!: BlogBrief[];
  displayColors: string[] = [];
  created() {
    this.displayColors = Object.getOwnPropertyNames(colors);
    // let count = Object.getOwnPropertyNames(colors).length;
    // while (this.displayColors.length < count) {
    //   let copiedColors: Colors = Copy(colors) as Colors;
    //   let displayColor = randomProperty(copiedColors) as string;
    //   delete copiedColors[displayColor];
    //   this.displayColors.push(displayColor);
    // }
  }
  rateChange(blog: BlogBrief) {
    // to be implemented
  }
}
</script>