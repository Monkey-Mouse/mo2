<template>
  <v-container class="fill-height">
    <v-breadcrumbs :items="items">
      <template v-slot:divider>
        <v-icon>mdi-forward</v-icon>
      </template>
    </v-breadcrumbs>
    <v-row>
      <v-col :key="i" v-for="(c, i) in cate">
        <v-hover v-slot="{ hover }">
          <v-card class="mx-auto" max-width="344">
            <v-card-text>
              <div>mo2 category</div>
              <p class="display-1 text--primary">{{ c.name }}</p>
            </v-card-text>
            <v-card-actions>
              <v-chip class="ma-2" color="secondary" label text-color="white">
                <v-icon left> mdi-label </v-icon>
                Category
              </v-chip>
            </v-card-actions>

            <v-expand-transition>
              <v-card
                v-if="hover"
                class="transition-fast-in-fast-out v-card--reveal"
                style="height: 100%"
                color="primary"
                @click="nextCate"
              >
                <v-card-text class="pb-0">
                  <p class="display-1 text--primary">Click to enter</p>
                  <!-- <p class="display-1 text--primary">Contains:</p>
                  <p class="display-1 text--secondary">
                    Sub Catagories + 64 articles
                  </p> -->
                </v-card-text>
              </v-card>
            </v-expand-transition>
          </v-card>
        </v-hover>
      </v-col>
    </v-row>
    <v-row justify="center">
      <a @click="$router.push('/')">Back to home</a>
    </v-row>
  </v-container>
</template>
<script lang="ts">
import { Category, User } from "@/models";
import { addQuery, GetCategories } from "@/utils";
import Vue from "vue";
import Component from "vue-class-component";
import { Prop } from "vue-property-decorator";
@Component({
  components: {},
})
export default class NotFound extends Vue {
  @Prop()
  user!: User;
  items = [
    {
      text: "Root",
      disabled: false,
      to: "?lv1=root",
    },
  ];
  cate: Category[] = [];
  lev = 1;
  created() {
    this.setCol("root");
    GetCategories(this.user.id).then((data) => {
      this.cate = data;
    });
  }
  setCol(name: string) {
    addQuery(this, `lv${this.lev++}`, name);
    addQuery(this, "cur", name);
  }
  nextLev(name: string) {
    this.items.push({
      text: name,
      disabled: false,
      to: `?lv${this.lev}=${name}`,
    });
    this.setCol(name);
  }
  nextCate() {
    this.nextLev("a");
  }
}
</script>
<style>
.v-card--reveal {
  bottom: 0;
  opacity: 1 !important;
  position: absolute;
  width: 100%;
}
</style>