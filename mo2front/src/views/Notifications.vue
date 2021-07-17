<template>
  <v-container class="fill-height">
    <v-row justify="center">
      <v-col align-self="center">
        <v-timeline align-top dense style="">
          <v-timeline-item
            v-for="(v, i) in datalist"
            :key="i"
            color="teal lighten-3"
            small
          >
            <v-row class="pt-1">
              <v-col cols="3">
                <time-ago
                  :refresh="60"
                  :datetime="v.create_time"
                  tooltip
                  long
                ></time-ago>
              </v-col>
              <v-col>
                <a
                  v-if="v.user !== undefined"
                  @click="$router.push(`/account/${v.operator_id}`)"
                  >{{ v.user.name }}</a
                ><span v-html="$sanitize(v.extra_message)"></span>
                <br />
                <avatar
                  v-if="v.user !== undefined"
                  class="mt-10"
                  :size="40"
                  :user="v.user"
                />
              </v-col>
            </v-row>
          </v-timeline-item>
        </v-timeline>
        <v-timeline align-top dense v-if="loading">
          <v-timeline-item v-for="i in 5" :key="i" color="teal lighten-3" small>
            <v-row class="pt-1">
              <v-col cols="3">
                <v-skeleton-loader type="text"></v-skeleton-loader>
              </v-col>
              <v-col>
                <v-skeleton-loader type="paragraph"></v-skeleton-loader>
              </v-col>
            </v-row>
          </v-timeline-item>
        </v-timeline>
        <nothing v-else-if="datalist.length === 0" />
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import { DisplayNotification, Notification, User } from "@/models";
import {
  GetNotifications,
  GetUserDatas,
  AutoLoader,
  ElmReachedBottom,
  AddMore,
} from "@/utils";
import Vue from "vue";
import Component from "vue-class-component";
import { Prop } from "vue-property-decorator";
import BlogSkeleton from "../components/BlogTimeLineSkeleton.vue";
import { TimeAgo } from "vue2-timeago";
import { Dictionary } from "node_modules/vue-router/types/router";
import Avatar from "../components/UserAvatar.vue";
import Nothing from "../components/NothingHere.vue";
@Component({
  components: {
    TimeAgo,
    Avatar,
    Nothing,
  },
})
export default class Notifications
  extends Vue
  implements AutoLoader<DisplayNotification>
{
  datalist: DisplayNotification[] = [];
  page = 0;
  pagesize = 10;
  loading: boolean = false;
  firstloading: boolean = false;
  nomore: boolean = false;
  created() {
    this.firstloading = true;
    this.loadMore();
  }
  loadMore() {
    this.loading = true;
    this.loadMoreData(this.page++, this.pagesize).then((data) => {
      AddMore(this, data);
    });
  }
  loadMoreData(page: number, pagesize: number) {
    return new Promise<DisplayNotification[]>((resolve, reject) => {
      GetNotifications({ page: page, pagesize: pagesize })
        .then((data) => {
          const dd = data as DisplayNotification[];
          const ids = data.map((v) => v.operator_id);
          GetUserDatas(ids)
            .then((d) => {
              const dic: Dictionary<User> = {};
              d.map((v) => (dic[v.id] = v as User));
              dd.forEach((v) => (v.user = dic[v.operator_id]));
              resolve(dd);
            })
            .catch((err) => reject(err));
        })
        .catch((err) => reject(err));
    });
  }
  ReachedButtom() {
    ElmReachedBottom(this, (q) => this.loadMoreData(q.page, q.pageSize));
  }
}
</script>
<style>
</style>
