<template lang="pug">
  header#header.app-header.navbar(role="menu")
    .navbar-header.bg-dark
      button.pull-right.visible-xs.dk()
      button.pull-right.visible-xs()
      a.navbar-brand.text-lt(href="/")
        img(src="../../images/logo.png")
        span.hidden-folded.m-l.xs Tania
    .collapse.pos-rlt.navbar-collapse.box-shadow.bg-white-only
      .nav.navbar-nav.hidden-xs
        a.btn.no-shadow.navbar-btn(:class="folded == true ? 'active': ''" @click="setFolded")
          i.fa.fa-dedent.text
          i.fa.fa-indent.text-active
      .nav.navbar-nav.hidden-xs
        .h3.text-lt.wrapper
          | Gracia Farm
          a.h6
            i.fa.fa-external-link
</template>

<script>
export default {
  name: 'AppHeaderComponent',
  props: {
    folded: { type: Boolean, default: false }
  },
  methods: {
    setFolded () {
      this.$emit('header-folded')
    }
  }
}
</script>
