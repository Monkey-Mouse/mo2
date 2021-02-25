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
        <v-card
          @click="gotoArticle(blog.id)"
          class="clickable mx-auto elevation-20"
          style="max-width: 400px"
        >
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
                  <div v-if="blog.userLoad" class="subtitle-1">
                    {{ blog.userName }}
                  </div>
                  <v-skeleton-loader v-else type="card-heading" />
                  <div class="subtitle-2">
                    {{ blog.entityInfo.createTime.substr(0, 10) }}
                  </div>
                </div>
              </v-card-title>
            </v-col>
            <v-col sm="6" cols="0">
              <v-img
                class="shrink ma-3"
                contain
                height="125px"
                :src="blog.cover"
                lazy-src="https://picsum.photos/id/11/100/60"
                style="flex-basis: 125px"
              >
                <template v-slot:placeholder>
                  <v-row
                    class="fill-height ma-0"
                    align="center"
                    justify="center"
                  >
                    <v-progress-circular
                      indeterminate
                      color="grey lighten-5"
                    ></v-progress-circular>
                  </v-row> </template
              ></v-img>
            </v-col>
          </v-row>
          <v-divider dark></v-divider>
          <v-card-actions
            v-on:click.prevent
            v-on:click.stop
            class="pa-4 unclickable"
          >
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
          <div class="py-4 text-break">
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
import { Prop, Watch } from "vue-property-decorator";
import { Colors } from "vuetify/es5/util/colors";
import { colors } from "vuetify/lib";
import { BlogBrief, DisplayBlogBrief } from "../models/index";
import { randomProperty, Copy, GetUserDatas } from "../utils/index";
@Component
export default class BlogTimeLineList extends Vue {
  @Prop()
  blogs!: DisplayBlogBrief[];
  prevlen = -1;
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
  mounted() {
    this.prevlen = this.blogs.length;
    const ids = this.blogs.map((v, i, a) => v.authorId);
    GetUserDatas(ids).then((data) => {
      const dic: any = {};
      for (let index = 0; index < data.length; index++) {
        const element = data[index];
        dic[element.id] = element.name;
      }
      for (let index = 0; index < this.blogs.length; index++) {
        this.blogs[index].userName = dic[this.blogs[index].authorId];
        this.blogs[index].userLoad = true;
      }
      this.$forceUpdate();
    });
  }
  @Watch("blogs")
  changeBlogs() {
    //TO DO: re fetch on loading new blogs
    // const ids = this.blogs.map((v, i, a) => v.authorId);
    // GetUserDatas(ids).then((data) => {
    //   const dic: any = {};
    //   for (let index = 0; index < data.length; index++) {
    //     const element = data[index];
    //     dic[element.id] = element.name;
    //   }
    //   for (let index = 0; index < this.blogs.length; index++) {
    //     this.blogs[index].userName = dic[this.blogs[index].authorId];
    //     this.blogs[index].userLoad = true;
    //   }
    // });
  }
  rateChange(blog: BlogBrief) {
    // to be implemented
  }
  gotoArticle(id: string) {
    this.$router.push("/article/" + id);
  }
}
</script>
<style>
.unclickable {
  cursor: default;
}
</style>