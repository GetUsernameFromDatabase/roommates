# Controllers

Controllers here use [swaggo/swag](https://github.com/swaggo/swag/tree/master?tab=readme-ov-file#api-operation) for API documentation

> Would like if it were able to automatically generate `@Route`.  
> They have a system in place to find data type, just find where the controller is used and look at it's path (and aggregate previous paths from group) then check if starts with base path. Only do that when `@Router` is `auto-gen` or something to signify "REPLACE_ME".

API controllers are more heavily seperated by files since the amount of documentation boilerplate required by swagger makes it harder to navigate
