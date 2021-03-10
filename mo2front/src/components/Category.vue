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
      <v-spacer />
      <v-col class="text-right">
        <v-btn small @click="addCate = true" fab color="primary">
          <v-icon> mdi-plus </v-icon>
        </v-btn></v-col
      >
    </v-row>
    <v-row v-if="cate.length > 0">
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
    <nothing v-else btnText="Create New" @click="addCate = true" />
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
  items = [
    {
      text: "Root",
      id: this.user.id,
    },
  ];
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
    const pos = this.items.indexOf(item);
    this.items = this.items.slice(0, pos + 1);
    const c = this.items[this.items.length - 1];
    this.parentId = c.id;
    this.loadData(c.id);
  }
  created() {
    this.setCol("root");
    this.loadData(this.user.id);
  }
  loadData(id: string) {
    GetCategories(id).then((data) => {
      this.cate = data;
    });
  }
  setCol(name: string) {
    addQuery(this, `lv${this.lev++}`, name);
    addQuery(this, "cur", name);
  }
  nextLev(c: Category) {
    this.parentId = c.id;
    this.loadData(c.id);
    const name = c.name;
    this.items.push({
      text: name,
      id: c.id,
    });
    this.setCol(name);
  }
  nextCate(c: Category) {
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