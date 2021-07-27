<template>
  <div>
    <v-timeline
      v-if="blogs && blogs.length > 0"
      :dense="this.$vuetify.breakpoint.smAndDown"
    >
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
                      :class="`${
                        displayColors[i % displayColors.length]
                      }--text`"
                      class="headline"
                      v-html="$sanitize(blog.title)"
                    ></div>
                    <a
                      v-on:click.prevent
                      v-on:click.stop
                      v-if="blog.userLoad"
                      @click="$router.push('/account/' + blog.authorId)"
                      class="subtitle-1"
                    >
                      <user-avatar :size="24" :user="blog.user" />
                      {{ blog.userName }}
                    </a>
                    <v-skeleton-loader v-else type="card-heading" />
                    <div class="subtitle-2">
                      <time-ago
                        :refresh="60"
                        :datetime="
                          blog.entityInfo
                            ? blog.entityInfo.updateTime
                            : blog['entityInfo.updateTime']
                        "
                        tooltip
                        long
                      ></time-ago>
                    </div>
                  </div>
                </v-card-title>
              </v-col>
              <v-col sm="6" cols="0">
                <v-img
                  class="shrink ma-3"
                  contain
                  height="125px"
                  :src="
                    (blog.cover
                      ? blog.cover
                      : '//cdn.mo2.leezeeyee.com/404.jpg') +
                    (blog.cover &&
                    blog.cover.indexOf('//cdn.mo2.leezeeyee.com') > 0
                      ? '~cover'
                      : '')
                  "
                  :lazy-src="
                    (blog.cover
                      ? blog.cover
                      : '//cdn.mo2.leezeeyee.com/404.jpg') +
                    (blog.cover &&
                    blog.cover.indexOf('//cdn.mo2.leezeeyee.com') > 0
                      ? '~thumb'
                      : '')
                  "
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
                        color="secondary"
                      ></v-progress-circular>
                    </v-row> </template
                ></v-img>
              </v-col>
            </v-row>
            <v-divider dark></v-divider>
            <v-card-actions class="pa-4">
              Rate
              <v-spacer></v-spacer>
              <span class="text--lighten-2 caption mr-2">
                ({{ blog.score_sum / blog.score_num }}) by
                {{ blog.score_num }} voter
              </span>
              <v-rating
                :value="blog.score_sum / blog.score_num"
                color="yellow accent-4"
                readonly
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
            class="justify-center row"
          >
            <div
              style="max-width: 400px; display: table; height: 200px"
              class="text-left col text-break"
            >
              <span
                style="
                  display: inline-block;
                  vertical-align: middle;
                  display: table-cell;
                "
                v-html="$sanitize(blog.description)"
              ></span>
            </div>
          </v-lazy>
        </template>
      </v-timeline-item>
    </v-timeline>
    <v-container v-else-if="showNothing">
      <v-row justify="center">
        <v-col class="text-center"
          ><v-icon size="128">mdi-clipboard-text-off-outline</v-icon></v-col
        >
      </v-row>
      <v-row justify="center">
        <v-col class="text-center"
          ><h1 class="text-h1">
            Nothing here yet :<span v-if="showI">(</span><span v-else>)</span>
          </h1>
        </v-col>
      </v-row>
      <v-row justify="center">
        <v-col class="text-center"
          ><v-icon size="128">mdi-arrow-down-bold</v-icon></v-col
        >
      </v-row>
      <v-row justify="center">
        <v-col class="text-center"
          ><v-btn outlined to="/edit" color="primary">
            Create your own article
          </v-btn>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import { Prop, Watch } from "vue-property-decorator";
import { Colors } from "vuetify/es5/util/colors";
import { colors } from "vuetify/lib";
import { BlogBrief, DisplayBlogBrief } from "../models/index";
import { randomProperty, Copy, GetUserDatas } from "../utils/index";
import UserAvatar from "./UserAvatar.vue";
import { TimeAgo } from "vue2-timeago";
@Component({
  components: {
    UserAvatar,
    TimeAgo,
  },
})
export default class BlogTimeLineList extends Vue {
  @Prop()
  blogs!: DisplayBlogBrief[];
  @Prop({ default: false })
  draft: boolean;
  @Prop({ default: true })
  showNothing: boolean;

  prevlen = -1;
  displayColors: string[] = [];
  showI = false;
  created() {
    this.displayColors = Object.getOwnPropertyNames(colors);
    setInterval(() => {
      this.showI = !this.showI;
    }, 500);
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
    if (ids.length !== 0) {
      GetUserDatas(ids).then((data) => {
        const dic: any = {};
        for (let index = 0; index < data.length; index++) {
          const element = data[index];
          dic[element.id] = element;
        }
        for (let index = 0; index < this.blogs.length; index++) {
          this.blogs[index].userName = dic[this.blogs[index].authorId].name;
          this.blogs[index].user = dic[this.blogs[index].authorId];
          this.blogs[index].userLoad = true;
          this.blogs[index].rate = 5;
        }
        this.$forceUpdate();
      });
    }
  }
  @Watch("blogs")
  changeBlogs() {
    if (this.blogs.length === this.prevlen) {
      // didn't add new blog, return
      return;
    }

    const ids = this.blogs.slice(this.prevlen).map((v, i, a) => v.authorId);
    GetUserDatas(ids).then((data) => {
      const dic: any = {};
      for (let index = 0; index < data.length; index++) {
        const element = data[index];
        dic[element.id] = element;
      }
      for (let index = this.prevlen; index < this.blogs.length; index++) {
        this.blogs[index].userName = dic[this.blogs[index].authorId].name;
        this.blogs[index].user = dic[this.blogs[index].authorId];
        this.blogs[index].userLoad = true;
        this.blogs[index].rate = 5;
      }
      this.$forceUpdate();
    });
  }
  rateChange(blog: BlogBrief) {
    // TODO to be implemented
  }
  gotoArticle(id: string) {
    this.$router.push(
      "/article/" + id + (this.draft ? `?draft=${this.draft}` : "")
    );
  }
}
</script>
<style>
.unclickable {
  cursor: default;
}
</style>