---
title: Pagination
description: Pagination
---

# Pagination

When a response returns many items or is paginated, the response contains a page channel instead of all the items. This ensures fast response times and the ability to just retrieve a required subset of items.

Responses that do return paginated results contain an item channel, that will emit pages of items.
