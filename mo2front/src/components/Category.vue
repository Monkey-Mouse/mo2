<template>
  <v-container class="fill-height">
    <MO2Dialog
      :validator="validator"
      :inputProps="inputProps"
      :show.sync="addCate"
      title="添加集合"
      confirmText="确认"
      :confirm="confirm"
    />
    <v-row>
      <v-col>
        <v-breadcrumbs :items="items">
          <template v-slot:divider>
            <v-icon>mdi-forward</v-icon>
          </template>
          <template v-slot:item="{ item }">
            <v-breadcrumbs-item
              :disabled="item.id === items[items.length - 1].id"
            >
              <a @click="gotoLevel(item)">{{ item.text.toUpperCase() }}</a>
            </v-breadcrumbs-item>
          </template>
        </v-breadcrumbs>
      </v-col>

      <v-btn
        v-if="own"
        class="mt-5 mr-3"
        small
        @click="addCate = true"
        fab
        color="primary"
      >
        <v-icon> mdi-plus </v-icon>
      </v-btn>
    </v-row>
    <v-row v-if="loading" justify="center">
      <v-sheet v-for="i in 4" :key="i">
        <v-skeleton-loader
          class="ma-4"
          min-width="200"
          type="card@2"
        ></v-skeleton-loader>
      </v-sheet>
    </v-row>
    <v-row v-else-if="cate.length > 0">
      <v-col :key="i" v-for="(c, i) in cate">
        <v-hover v-slot="{ hover }">
          <v-card elevation="10" class="mx-auto" max-width="344">
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
                elevation="10"
                v-if="hover"
                class="transition-fast-in-fast-out v-card--reveal"
                style="height: 100%"
                color="primary"
                @click="nextCate(c)"
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
    <nothing
      v-else
      :btnText="own ? 'Create New' : ''"
      @click="addCate = true"
    />
    <v-row justify="center">
      <v-divider></v-divider>
    </v-row>
  </v-container>
</template>
<script lang="ts">
import { Category, InputProp, User } from "@/models";
import { addQuery, GetCategories, GetErrorMsg, UpsertCate } from "@/utils";
import Vue from "vue";
import Component from "vue-class-component";
import { Prop } from "vue-property-decorator";
import { required } from "vuelidate/lib/validators";
import MO2Dialog from "./MO2Dialog.vue";
import Nothing from "./NothingHere.vue";
@Component({
  components: { MO2Dialog, Nothing },
})
export default class Mo2Category extends Vue {
  @Prop()
  user!: User;
  @Prop()
  own: boolean;
  items: { text: string; id: string }[] = [];
  cate: Category[] = [];
  lev = 1;
  parentId = "";
  validator = {
    name: {
      required: required,
    },
  };
  inputProps: { [name: string]: InputProp } = {
    name: {
      errorMsg: {
        required: "集合名不可为空",
      },
      label: "Name",
      default: "",
      icon: "mdi-folder",
      col: 12,
      type: "text",
    },
  };
  addCate = false;
  loading = true;
  async confirm({ name }: { name: string }) {
    try {
      const data = await UpsertCate({
        name: name,
        parent_id: this.parentId === "" ? this.user.id : this.parentId,
      });
      this.cate.push(data);
      return { err: "", pass: true };
    } catch (error) {
      return { err: GetErrorMsg(error), pass: false };
    }
  }
  gotoLevel(item: { text: string; id: string }) {
    this.loading = true;
    const pos = this.items.indexOf(item);
    this.lev = pos + 1;
    this.items = this.items.slice(0, pos + 1);
    const c = this.items[pos];
    this.parentId = c.id;
    addQuery(this, "cur", c.id);
    this.loadData(c.id);
  }
  created() {
    this.items.push({
      text: "Root",
      id: this.user.id,
    });
    const urlParams = new URLSearchParams(window.location.search);
    const cur = urlParams.get("cur");
    if (cur) {
      this.items = [];
      for (
        let index = 0;
        index < Object.keys(this.$route.query).length;
        index++
      ) {
        const element = this.$route.query[`lvid${index + 1}`] as string;
        const name = this.$route.query[`lv${index + 1}`] as string;
        this.items.push({ text: name, id: element });
        if (element === cur) {
          break;
        }
      }
      this.lev = this.items.length + 1;
      // if (this.items[0].id != this.user.id) {
      //   this.items.unshift({
      //     text: "Root",
      //     id: this.user.id,
      //   });
      // }
      this.loadData(cur);
      return;
    }
    this.setCol("root", this.user.id);
    this.loadData(this.user.id);
  }
  loadData(id: string) {
    this.loading = true;
    GetCategories(id).then((data) => {
      this.cate = data;
      this.loading = false;
    });
  }
  setCol(name: string, id: string) {
    // addQuery(this, `lvid1`, this.user.id);
    addQuery(this, `lv${this.lev}`, name);
    addQuery(this, `lvid${this.lev}`, id);
    this.lev++;
    addQuery(this, "cur", id);
  }
  nextLev(c: Category) {
    this.parentId = c.id;
    this.loadData(c.id);
    const name = c.name;
    this.items.push({
      text: name,
      id: c.id,
    });
    this.setCol(name, c.id);
  }
  nextCate(c: Category) {
    this.loading = true;
    this.nextLev(c);
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