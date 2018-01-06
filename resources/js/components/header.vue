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
        .h3.text-lt.wrapper
          | Gracia Farm
          a.h6
            i.fa.fa-external-link
</template>

<script>
export default {
  name: 'AppHeaderComponent',
}
</script>
