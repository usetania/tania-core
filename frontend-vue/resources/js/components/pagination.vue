<template lang="pug">
  div
    .panel-footer(v-if="pages > 1")
      .text-center
        ul.pagination.pagination-sm.m-t-none.m-b-none
          li(v-bind:class="{ disabled: currentPage == 1 }" v-on:click="reload")
            router-link(:to="{ path: path, query: { page: (currentPage - 1) }}")
              i.fa.fa-chevron-left
          li(v-if="pages <= tabs" v-for="pageNumber in pages" v-on:click="reload")
            router-link(:to="{ path: path, query: { page: pageNumber }}" v-bind:class="{ active: currentPage == pageNumber }") {{ pageNumber }}
          li(v-if="pages > tabs && currentPage > mid && (currentPage + span) <= pages" v-for="pageNumber in tabs" v-on:click="reload")
            router-link(:to="{ path: path, query: { page: (currentPage - mid + pageNumber) }}" v-bind:class="{ active: currentPage == (currentPage - mid + pageNumber) }") {{ currentPage - mid + pageNumber }}
          li(v-if="pages > tabs && currentPage > mid && (currentPage + span) > pages && currentPage < pages " v-for="pageNumber in tabs" v-on:click="reload")
            router-link(:to="{ path: path, query: { page: (currentPage - (mid + (pages - currentPage)) + pageNumber)  }}" v-bind:class="{ active: currentPage == (currentPage - (mid + (pages - currentPage)) + pageNumber) }") {{ currentPage - (mid + (pages - currentPage)) + pageNumber }}
          li(v-if="pages > tabs && currentPage > mid && currentPage == pages" v-for="pageNumber in tabs" v-on:click="reload")
            router-link(:to="{ path: path, query: { page: (currentPage - (tabs - pageNumber)) }}" v-bind:class="{ active: currentPage == (currentPage - (tabs - pageNumber)) }") {{ (currentPage - (tabs - pageNumber)) }}
          li(v-if="pages > tabs && currentPage <= mid" v-for="pageNumber in tabs" v-on:click="reload")
            router-link(:to="{ path: path, query: { page: pageNumber }}" v-bind:class="{ active: currentPage == pageNumber }") {{ pageNumber }}
          li(v-bind:class="{ disabled: currentPage == pages }" v-on:click="reload")
            router-link(:to="{ path: path, query: { page: (currentPage + 1) }}")
              i.fa.fa-chevron-right
</template>

<script>
export default {
  name: 'Pagination',
  data () {
    return {
      currentPage: 1,
      tabs: 5,
      mid: 3,
      span: 2,
      path: '',
    }
  },
  methods: {
    reload () {
      this.$emit('reload')
      this.currentPage = parseInt(this.$route.query.page)
    }
  },
  mounted () {
    this.path = this.$route.path
    if (typeof this.$route.query.page != "undefined") {
      if (parseInt(this.$route.query.page) <= this.pages && this.$route.query.page > 0) {
        this.currentPage = parseInt(this.$route.query.page)
      }
      else {
        this.$router.replace(this.$route.path)
      }
    }
    this.$watch('current', current => {
      this.currentPage = this.current
    }, {})
  },
  props: ['pages', 'current'],
}
</script>
