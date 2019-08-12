;(function ($) {
    var viewerData = function () {
        var data = null
        var orderBy = ''
        var orderType = 'asc'
        var page = '1'
        var vid = ''
        var filter=[]
        var iofunc={}
        var makeIOfunc=function(script){
            return eval('(function(js){if(js===""||js===undefined){return {in:function(data){return data}, out:function(data){return data}}}else{return js}})('+script+')')
        }
        var initIO=function(data){
            for(var key in data){
                var script=atob(data[key])
                iofunc[key]=makeIOfunc(script)
            }
        }
        var io=function(name){
            return iofunc[name]
        }
        var getInfo = function () {
            if (data == null) {
                return null
            }
            return data['info']
        }
        var getMap = function () {
            if (data == null) {
                return null
            }
            return data['info']['map']
        }
        var getFields = function () {
            if (data == null) {
                return null
            }
            return data['info']['fields']
        }
        var primaryKey = function () {
            return data['info']['primaryKey']
        }
        var primaryName = function () {
            return data['info']['primaryName']
        }
        var getVid = function () {
            return vid
        }
        var get = function (viewerID,ft) {
            if (viewerID !== vid) {
                orderBy = ""
                orderType = 'asc'
            }
            if(ft===undefined){
                ft=[]
            }
            $.ajax({
                url: '/api/viewerdata?' + 'viewer-id=' + viewerID + '&limit=20+' + '&order=' + orderBy + '&order-type=' + orderType + '&page=1'+'&filter='+JSON.stringify(ft),
                method: 'GET',
                async: 'true',
                cache: 'false',
                success: function (d) {
                    data = d
                    if(vid != viewerID){
                        iofunc={}
                        initIO(data['info']['script'])
                        vid = viewerID
                    }
                    page = data['pages']['current']
                    filter=ft
                    if (orderBy === "") {
                        orderBy = data['info']['primaryKey']
                    }
                    drawOperator()
                    drawPagebar()
                    drawTable()
                    if(filter.length!==0){
                        $('#clear-filter').removeClass('invisible')
                    }else{
                        $('#clear-filter').addClass('invisible')
                    }
                },
                error: function (err) {
                    message(err.responseJSON.msg, 'danger')
                }
            })
        }
        var getPage = function (p) {
            $.ajax({
                url: '/api/viewerdata?' + 'viewer-id=' + vid + '&limit=20+' + '&order=' + orderBy + '&order-type=' + orderType + '&page=' + p + '&filter='+JSON.stringify(filter),
                method: 'GET',
                async: 'true',
                cache: 'false',
                success: function (d) {
                    data = d
                    page = data['pages']['current']
                    drawOperator()
                    drawPagebar()
                    drawTable()
                    if(filter.length!==0){
                        $('#clear-filter').removeClass('invisible')
                    }else{
                        $('#clear-filter').addClass('invisible')
                    }
                },
                error: function (err) {
                    message(err.responseJSON.msg, 'danger')
                }
            })
        }
        var clearFilter=function(){
            filter=[]
            var searcher = $('#viewer-data-searcher form')
            searcher.attr('vid','')
            get(vid)
        }
        var flush=function(){
            getPage(page)
        }
        var drawOperator = function(){
            var operaterContainer=$('.table-operator')
            operaterContainer.html('')
            var deleter=$('<button type="button" class="btn btn-danger" id="data-deleter">删除</button>')
            var adder=$('<button type="button" class="btn btn-primary" id="data-adder">添加</button>')
            adder.css('margin-left','5px')
            var searcher=$('<button class="btn btn-primary" id="data-searcher">搜索</button>')
            searcher.css('margin-left','5px')
            var filterClearer=$('<button class="btn btn-primary invisible" id="clear-filter">清除过滤</button>')
            filterClearer.css('margin-left','5px')
            var pageList=$('<ul class="float-sm-right pagination pagination-sm" id="page-bar"></ul>')
            pageList.css('margin-left','5px')
            adder.click(function () {
                if($.viewerData.getVid()==='1'){
                    var editer=$('#menu-item-editer')
                    editer.find('.modal-title').html('添加菜单项')
                    editer.find('.ensure').attr('action','add')
                    editer.modal()
                }else if($.viewerData.getVid()==='2'){
                    var editer=$('#viewer-item-editer')
                    editer.find('.modal-title').html('添加视图项')
                    editer.find('.ensure').attr('action','add')
                    editer.modal()
                }else{
                    var fields = $.viewerData.getFields()
                    if (fields == null) {
                        return
                    }
                    var data = []
                    for (var i = 0; i < fields.length; i++) {
                        data.push(fields[i]['name'])
                    }
                    $.itemEditer.creater(data)
                }
            })
            deleter.click(function () {
                var data = []
                var checker = $('.row-checker')
                for (var i = 0; i < checker.length; i++) {
                    var c = $(checker[i])
                    if (c.is(':checked')) {
                        var row = c.parent().parent()
                        var primaryTd = row.children('td[col="' + $.viewerData.primaryName() + '"]')
                        var primary = primaryTd.children('span')
                        data.push(($.viewerData.io($.viewerData.primaryName()).out(primary.html())).toString())
                    }
                }
                $.ajax({
                    url: '/api/viewerdata?' + 'param=' + JSON.stringify(data) + '&viewer-id=' + $.viewerData.getVid(),
                    method: 'DELETE',
                    sync: 'false',
                    cache: 'false',
                    success: function (data) {
                        message('删除成功', 'success')
                        $.viewerData.flush()
                    },
                    error: function (err) {
                        message(err.responseJSON.msg, 'danger')
                    }
                })
            })
            searcher.click(function(){
                var fields = $.viewerData.getInfo()
                if (fields == null) {
                    return
                }
                $.itemEditer.searcher(fields)
            })
            filterClearer.click(function(){
                $.viewerData.clearFilter()
            })
            operaterContainer.append(deleter)
            operaterContainer.append(adder)
            operaterContainer.append(searcher)
            operaterContainer.append(filterClearer)
            operaterContainer.append(pageList)
        }
        var drawPagebar = function () {
            if (data == null) {
                return
            }
            var pageBar = $('#page-bar')
            pageBar.html('')
            var pageInfo = data['pages']
            var page
            if (pageInfo['first'] != null) {
                page = $('<li class="page-item"><a class="page-link page-num" href="#" page="1">首页</a></li>')
                pageBar.append(page)
            }
            if (pageInfo['prev'] != null) {
                page = $('<li class="page-item"><a class="page-link page-num" href="#" page="' + pageInfo['prev'] + '">' + '<上一页</a></li>')
                pageBar.append(page)
            }
            var pages = pageInfo['pages']
            for (var i in pages) {
                page = $('<li class="page-item"><a class="page-link page-num" href="#" page="' + pages[i] + '">' + pages[i] + '</a></li>')
                if (pages[i] === pageInfo['current']) {
                    page.addClass('active')
                }
                pageBar.append(page)
            }
            if (pageInfo['total'] != null) {
                page = $('<li class="page-item"><a class="page-link page-num" href="#" page="' + pageInfo['total'] + '">' + '...' + pageInfo['total'] + '</a></li>')
                pageBar.append(page)
            }
            if (pageInfo['next'] != null) {
                page = $('<li class="page-item"><a class="page-link page-num" href="#" page="' + pageInfo['next'] + '">' + '下一页></a></li>')
                pageBar.append(page)
            }
            if (pageInfo['last'] != null) {
                page = $('<li class="page-item"><a class="page-link page-num" href="#" page="' + pageInfo['last'] + '">' + '尾页</a></li>')
                pageBar.append(page)
            }
            var pelem = $('#page-bar a')
            pelem.click(function () {
                var page = $(this).attr('page')
                getPage(page)
            })
        }
        var drawTable = function () {
            if (data == null) {
                return
            }
            var table = $('#main-table')
            table.html('')
            var thead = $('<thead  class="bg-warning"></thead>')
            var row = $('<tr></tr>')
            var col = $('<td><input type="checkbox" id="table-checkbox"></td>')
            row.append(col)
            var fields = data.info.fields
            for (var index in fields) {
                var order = ""
                if (fields[index]['order'] === true) {
                    var orderFlag = '▲'
                    if (orderType === 'desc') {
                        orderFlag = '▼'
                    }
                    if (fields[index]['field'] === orderBy) {
                        order = '<span class="text-primary">' + orderFlag + '</span>'
                    } else {
                        order = orderFlag
                    }
                }
                col = $('<td name="' + fields[index]['field'] + '">' + fields[index]['name'] + order + '</td>')
                if (fields[index]['order'] === true) {
                    col.addClass('can-sort')
                    col.css("cursor", "pointer")
                    col.click(function () {
                        var orderKey = $(this).attr('name')
                        if (orderKey === orderBy) {
                            if (orderType === 'asc') {
                                orderType = 'desc'
                            } else {
                                orderType = 'asc'
                            }
                        } else {
                            orderBy = orderKey
                        }
                        flush()
                    })
                }
                row.append(col)
            }
            if(vid ==='1'||vid==='2'){
                col=$('<td></td>')
                row.append(col)
            }
            thead.append(row)
            table.append(thead)
            var tbody = $('<tbody></tbody>')
            var content = data.data
            for (var i = 0; i < content.length; i++) {
                row = $('<tr></tr>')
                col = $('<td><input class="row-checker" type="checkbox"></td>')
                row.append(col)
                for (var index in fields) {
                    col = $('<td col="' + fields[index]['name'] + '">' + '<span class="td-content">' + $.viewerData.io(fields[index]['name']).in(content[i][index]) + '</span>' + '</td>')
                    var span = col.children('span')
                    span.css('cursor', 'pointer')
                    span.click(function () {
                        var td = $(this).parent()
                        var row = td.parent()
                        var col = td.attr('col')
                        var v = $(this).html()
                        var primaryTd = row.children('td[col="' + $.viewerData.primaryName() + '"]')
                        var primary = primaryTd.children('span').html()
                        var data = {'key': col, 'value': v, 'primary': primary}
                        $.itemEditer.editer(data)
                        $('#viewer-data-editer').modal()
                    })
                    row.append(col)
                }
                if(vid ==='1'||vid==='2'){
                    var button=$('<button type="button" class="btn btn-primary btn-sm">编辑</button>')
                    if(vid==='1'){
                        button.click(function(){
                            var editer=$('#menu-item-editer')
                            var form = editer.find('form')
                            var row=$(this).parent().parent()
                            var name=row.children('td[col="菜单名"]').children('span').html()
                            var allow=row.children('td[col="访问者"]').children('span').html()
                            var nameInput= form.children('input[name="item-name"]')
                            var allowInput= form.children('input[name="allow"]')
                            var id=row.children('td[col="菜单ID"]').children('span').html()
                            nameInput.val(name)
                            allowInput.val(allow)
                            editer.find('.modal-title').html('编辑菜单项')
                            editer.find('.ensure').attr('action','update')
                            editer.find('.ensure').attr('menu-id',id)
                            editer.modal()
                        })
                    }else if(vid==='2'){
                        button.click(function(){
                            var editer=$('#viewer-item-editer')
                            var baseform = $('#viewer-base-info')
                            var row=$(this).parent().parent()
                            var nameInput= baseform.children('input[name="item-name"]')
                            var allowInput= baseform.children('input[name="allow"]')
                            var sqlform= $('#viewer-item-editer .sql-config')
                            var dbAddrInput=sqlform.find('input[name="address"]')
                            var dbPassInput=sqlform.find('input[name="password"]')
                            var dbNameInput=sqlform.find('input[name="database"]')
                            var dbTableInput=sqlform.find('input[name="table"]')
                            var dbUserInput=sqlform.find('input[name="user"]')
                            var name=row.children('td[col="视图名"]').children('span').html()
                            var allow=row.children('td[col="访问者"]').children('span').html()
                            var dbAddr=row.children('td[col="数据库地址"]').children('span').html()
                            var dbPass=row.children('td[col="数据库密码"]').children('span').html()
                            var dbName=row.children('td[col="数据库名"]').children('span').html()
                            var dbTable=row.children('td[col="表名"]').children('span').html()
                            var dbUser=row.children('td[col="数据库用户"]').children('span').html()
                            nameInput.val(name)
                            allowInput.val(allow)
                            dbAddrInput.val(dbAddr)
                            dbPassInput.val(dbPass)
                            dbNameInput.val(dbName)
                            dbTableInput.val(dbTable)
                            dbUserInput.val(dbUser)
                            var data=$.func.getFormData($('#viewer-item-editer .sql-config'))
                            $.ajax({
                                url:'/api/databaseField?'+data,
                                method:'GET',
                                async:false,
                                cache:false,
                                success:function(data){
                                    $.dbFields.clear()
                                    $.dbFields.set(data)
                                    $('#sql-check-result').removeClass('text-danger')
                                    $('#sql-check-result').addClass('text-success')
                                    $('#sql-check-result').text('检查数据库成功')
                                },
                                error:function(err){
                                    $.dbFields.clear()
                                    $('#sql-check-result').removeClass('text-success')
                                    $('#sql-check-result').addClass('text-danger')
                                    $('#sql-check-result').text(err.responseJSON.msg)
                                },
                            })
                            try{
                                var fieldsInfo=JSON.parse(row.children('td[col="字段信息"]').children('span').html())
                                var info=fieldsInfo['fields']
                                var fields=$.dbFields.get()
                                var option=''
                                for(var i=0;i<fields.length;i++){
                                    option=option+'<option>'+fields[i]['name']+'</option>'
                                }
                                var container=$('#viewer-item-editer .viewer-fields')
                                for(var i=0;i<info.length;i++){
                                    var group=$('<div class="form-group form-inline viewer-field alert"><label>字段名：</label></div>')
                                    var nameInput=$('<input type="text" name="field-name">')
                                    nameInput.val(info[i]['name'])
                                    var select=$('<select class="custom-select custom-select-sm">'+option+'</select>')
                                    select.val(info[i]['field'])
                                    var canOrder=$('<input type="checkbox" class="form-check-input" name="can-order" value="">')
                                    if(info[i]['order']===true){
                                        canOrder.prop("checked", true)
                                    }
                                    var primary=$('<input type="radio" name="primary-key">')
                                    if(info[i]['name']===fieldsInfo['primaryName']){
                                        primary.prop("checked",true)
                                    }
                                    var scriptID = randomString(8, '0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ')
                                    var scriptBtn=$('<button type="button" class="btn btn-primary btn-sm" name="script-creater" script-id="'+scriptID+'" data-toggle="modal" data-target="#viewer-data-script-editer">脚本</button>')
                                    scriptBtn.click(function(){
                                        var scriptID=$(this).attr('script-id')
                                        var textArea=$('#viewer-data-script-editer textarea')
                                        textArea.attr('script-id',scriptID)
                                        var script=$('#'+scriptID).val()
                                        textArea.val(script)
                                    })
                                    var scriptText=$('<textarea class="d-none" name="script-container" id="'+scriptID+'"></textarea>')
                                    var script=atob(fieldsInfo['script'][info[i]['name']])
                                    scriptText.val(script)
                                    if(script!==''){
                                        scriptBtn.removeClass('btn-primary')
                                        scriptBtn.addClass('btn-danger')
                                    }
                                    var dissmiss=$('<button type="button" class="close" data-dismiss="alert">&times;</button>')
                                    group.append(nameInput)
                                    group.append(select)
                                    group.append(canOrder)
                                    group.append('可排序')
                                    group.append(primary)
                                    group.append('主键')
                                    group.append(scriptBtn)
                                    group.append(scriptText)
                                    group.append(dissmiss)
                                    container.append(group)
                                }
                                var id=row.children('td[col="视图ID"]').children('span').html()
                                editer.find('.modal-title').html('编辑视图项')
                                editer.find('.ensure').attr('action','update')
                                editer.find('.ensure').attr('viewer-id',id)
                            }catch(error){
                                message(error.message,'danger')
                            }finally {
                                editer.modal()
                            }
                        })
                    }
                    col=$('<td></td>')
                    col.append(button)
                    row.append(col)
                }
                tbody.append(row)
            }
            table.append(tbody)
            $('#table-checkbox').change(function () {
                if ($(this).is(':checked')) {
                    $('#main-table input[type="checkbox"]').prop("checked", true)
                } else {
                    $('#main-table input[type="checkbox"]').prop("checked", false)
                }
            })
        }
        return {
            get: get,
            getPage: getPage,
            drawOperator:drawOperator,
            drawPagebar: drawPagebar,
            drawTable: drawTable,
            primaryKey: primaryKey,
            primaryName: primaryName,
            getVid: getVid,
            getMap: getMap,
            getFields: getFields,
            getInfo: getInfo,
            flush:flush,
            clearFilter:clearFilter,
            io:io,
        }
    }
    $.extend({
        viewerData: viewerData(),
    })
    /*****************************************************************************
     * 编辑器
     */
    var itemEditer = function () {
        var editer = function (data) {
            var header = $('#viewer-data-editer .header-title')
            var body = $('#viewer-data-editer .modal-body')
            header.html(data['key'])
            var text = $('<textarea class="form-control" rows="1" field="' + data['key'] + '"' + 'primary="' + data['primary'] + '"' + '>' + data['value'] + '</textarea>')
            text.on('input', function () {
                this.style.height = 'auto'
                this.style.height = (this.scrollHeight) + 'px'
            })
            body.append(text)
        }
        var creater = function (data) {
            var container = $('#viewer-data-creater form')
            var info = $.viewerData.getInfo()
            for (var i = 0; i < data.length; i++) {
                var def = info['default'][data[i]]
                var nullable = info['nullable'][data[i]]
                var extra = info['extra'][data[i]]
                var deftext = '   <small class="text-secondary">默认(' + def + ')' + '</small>'
                if (extra === 'auto_increment') {
                    deftext = '   <small>默认(自增)</small>'
                }
                var group = $('<div class="form-group"></div>')
                var label = $('<label>' + data[i] + deftext + '</label>')
                var text = $('<textarea class="form-control" rows="1" field="' + data[i] + '"' + '>' + '</textarea>')
                if (def == null && nullable === 'NO' && extra !== 'auto_increment') {
                    text.addClass('is-invalid')
                }
                text.on('input', function () {
                    this.style.height = 'auto'
                    this.style.height = (this.scrollHeight) + 'px'
                })
                group.append(label)
                group.append(text)
                if (extra === 'auto_increment') {
                    var fs = $('<fieldset disabled></fieldset>')
                    group = fs.append(group)
                }
                container.append(group)
            }
            $('#viewer-data-creater').modal()
        }
        var searcher=function(data){
            var fields=data['fields']
            var container = $('#viewer-data-searcher form')
            var vid=container.attr('vid')
            if(vid!==$.viewerData.getVid()) {
                container.html('')
                container.attr('vid',$.viewerData.getVid())
                for (var i = 0; i < fields.length; i++) {
                    var name = fields[i]['name']
                    var field = fields[i]['field']
                    var group = $('<div class="form-group"></div>')
                    var label = $('<label>' + name + '</label>')
                    var selector = $('<select><option>=</option><option>&gt</option><option>&lt</option></select>')
                    var text = $('<textarea class="form-control" rows="1" name="' + name + '" field="' + field + '"' + '>' + '</textarea>')
                    text.on('input', function () {
                        this.style.height = 'auto'
                        this.style.height = (this.scrollHeight) + 'px'
                    })
                    group.append(label)
                    group.append(selector)
                    group.append(text)
                    container.append(group)
                }
            }
            $('#viewer-data-searcher').modal()
        }
        return {
            editer: editer,
            creater: creater,
            searcher:searcher,
        }
    }
    $.extend({
        itemEditer: itemEditer(),
    })
})(jQuery)
;(function ($) {
    var viewerFields=function(container){
        var fieldsInfo=$.dbFields.get()
        var info={'primaryName':null,'primaryKey':null,'fields':[],'map':{},'default':{},'nullable':{},'extra':{},'script':{}}
        var containers=container.children('.viewer-field')
        for(var i=0;i<containers.length;i++){
            var name=$(containers[i]).children('input[name="field-name"]').val()
            if(name===''){
                continue
            }
            if((info['map']).hasOwnProperty(name)){
                return [null,'字段名重复']
            }
            var field=$(containers[i]).children('select').val()
            var order=$(containers[i]).children('input[name="can-order"]').is(':checked')
            var primary=$(containers[i]).children('input:radio[name="primary-key"]').is(':checked')
            if(primary){
                info.primaryName=name
                info.primaryKey=field
            }
            info['fields'].push({'name':name,'field':field,'order':order})
            info['map'][name]=field
            for(var j=0;j<fieldsInfo.length;j++){
                if(fieldsInfo[j]['name']===field){
                    info['default'][name]=fieldsInfo[j]['default']
                    info['nullable'][name]=fieldsInfo[j]['nullable']
                    info['extra'][name]=fieldsInfo[j]['extra']
                    break
                }
            }
            var scriptContainer=$(containers[i]).children('textarea[name="script-container"]')
            var script = scriptContainer.val()
            script=btoa(script)
            info['script'][name]=script
        }
        return [info,null]
    }
    var dbFields=function(){
        var fields
        return {
            set:function(f){fields=f},
            clear:function(){fields=null},
            get:function(){return fields},
            viewerFields:viewerFields,
        }
    }
    $.extend({
        dbFields:dbFields(),
    })
    var getFormData=function(form){
        var data=form.serialize()
        return decodeURIComponent(data,true)
    }
    func={
        getFormData:getFormData,
    }
    $.extend({
        func:func,
    })
    var getMenuItem=function(parentID){
        var result
        $.ajax({
            url:'/api/menuitem?parent='+parentID,
            method:'GET',
            async:false,
            cache:false,
            success:function(data){
                result=[true,data]
            },
            error:function(err){
                result=[false,err.responseJSON]
            },
        })
        return result
    }
    var getViewerItem=function(parentID){
        var result
        $.ajax({
            url:'/api/vieweritem?parent='+parentID,
            method:'GET',
            async:false,
            cache:false,
            success:function(data){
                result=[true,data]
            },
            error:function(err){
                result=[false,err.responseJSON]
            },
        })
        return result
    }
    api={
        getMenuItem:getMenuItem,
        getViewerItem:getViewerItem,
    }
    $.extend({
        api: api
    })
    var createSelector = function (container, parentID) {
        var result=$.api.getMenuItem(parentID)
        if (result[0]){
            data=result[1]
            if (data.items == null) {
                return
            }
            var option = '<option>选择父节点</option>\n'
            for (var i = 0; i < data.items.length; i++) {
                option = option + '<option item-id="' + data.items[i].id + '">' + data.items[i].name + '</option>\n'
            }
            var e = $('<div class="alert alert-dismissable alert-danger"><button type="button" class="close" data-dismiss="alert">&times;</button>' +
                '<select class="custom-select parent-selector" onchange="$.selector.onParentSelected(this)" isselected="false">\n' +
                option +
                '</select></div>')
            container.append(e)
        }else{
            data=result[1]
            message(data.msg, 'danger')
        }
    }
    var getSelectedParentID = function (selectors) {
        var parentID = 0
        var parentName=""
        if (selectors.length > 0) {
            var selector = selectors.last()
            var option = $(selector.get(0)[selector.get(0).selectedIndex])
            parentID = option.attr('item-id')
            parentName=selector.val()
        }
        return [parentID,parentName]
    }
    var onParentSelected = function (e) {
        var selector = $(e)
        if (selector.get(0).selectedIndex > 0) {
            if (selector.parent().hasClass('alert-danger')) {
                selector.parent().removeClass('alert-danger')
                selector.parent().addClass('alert-success')
            }
            if (selector.attr('isselected') == 'false') {
                selector.attr('isselected', true)
            }
        } else {
            if (selector.parent().hasClass('alert-success')) {
                selector.parent().removeClass('alert-success')
                selector.parent().addClass('alert-danger')
            }
            if (selector.attr('isselected') == 'true') {
                selector.attr('isselected', false)
            }
        }
    }
    var hasUnselectedSelector = function (selectors) {
        for (var i = 0; i < selectors.length; i++) {
            if ($(selectors[i]).attr('isselected') == 'false') {
                return true
            }
        }
        return false
    }
    var selector = {
        createSelector: createSelector,
        getSelectedParentID: getSelectedParentID,
        onParentSelected: onParentSelected,
        hasUnselectedSelector: hasUnselectedSelector,
    }
    $.extend({
        selector: selector
    })
})(jQuery)