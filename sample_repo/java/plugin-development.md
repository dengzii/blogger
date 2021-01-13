插件作为提升工作效率的工具, 很大程度的加快了我们的开发速度, 给我们的工作带来了极大的便利. 在编码的时候, 可曾设想过如果 IDE 有 xx 功能
该多好, 相信大部分人都有过这种想法. 这篇文章记录了我学习开发 IntelliJ IDEA 插件的过程.

**开发环境**

- 系统: Windows 10
- 工具: IntelliJ IDEA 2019.2.1 Community Edition
- SDK: Java 8, Kotlin 1.3.41
- AndroidStudio: Android Studio 3.5

**官方文档**

> http://www.jetbrains.org/intellij/sdk/docs/reference_guide

**我开发的插件**

- [WiFiADB](https://github.com/MrDenua/WiFiADB) : WiFi 连接 Android 手机 ABD

- [FindViewGenerator](https://github.com/MrDenua/FindViewGenerator) : 自动生成 findViewById 代码

### 插件的大致分类

#### 语言支持

For example, Gradle, Scala, Groovy, 这些插件 IDEA, AndroidStudio 都是自带有的, 所以我们才能在编写这些语言代码的时候有语法高亮检测.

语言插件主要的一些功能是 文件类型识别, 语法高亮检测, 格式化, 语法提示等等.

#### 框架集成

AndroidStudio 就是一个例子, 他集成了 AndroidSDK 的一系列功能,  比如, 资源文件识别组织与提示, 集成 Gradle, Debug, ADB, 打包APK, 让我们可以更好的开发 Android 应用. 类似的插件还有 Gradle, Maven, Spring Plugin等.

框架集成插件的主要功能是某种语言特定代码的识别, 直接访问特定框架的功能

#### 工具集成

例如, 翻译工具, Markdown View, WiFi ADB等.

#### UI增强

对 IDE 的主题颜色样式做一些更改, 比如 MaterialTheme.

## 创建一个项目

#### 两种工具开发插件项目

**Gradle**

本着紧随时代潮流的想法, 刚开始我是用 Gradle 构建项目的, 但是, 我在 IDEA Community 2019.2.1 版本中, 无论如何都无法成功, Gradle 提示找不到 PsiJavaFile 这些类, 但项目中是可以引用的. 我尝试了换 Gradle版本, 5.4, 4.5.1 两个版本, 把 jar 包移到 libs 中依赖, 均以失败告终, 如果有人可以编译运行, 请千万要告诉我.

**DevKit**

创建项目:

    New Project => IntelliJ Platform Plugin => Input Project Name => Finish

配置项目:

    File => Project Structure
				Project => Project SDK => IntelliJ IDEA Community Edition IC-xxx
				Module => Select your module => Tab Dependencies => Module SDK => Project


创建完毕后, 你的目录结构应如下

	resource/
		META-INF/
			plugin.xml	// plugin config file
	src/	// source code directory

#### plugin.xml

```html
<idea-plugin url="https://www.your_plugin_home_page.com">	

    <name>Your plugin name</name>

    <id>com.your_domain.plugin_name</id>

    <depends>com.intellij.modules.all</depends>
    <!-- kotlin support -->		
    <depends>org.jetbrains.kotlin</depends>
    
    <description>Your will see it at plugin download page</description>

    <change-notes>What's update</change-notes>

    <version>1.0.0</version>
</idea-plugin>
```

如果你的插件需要支持 kotlin, 则必须添加这个依赖

	<depends>org.jetbrains.kotlin</depends>

## 准备工作

#### 线程规则

在 IntelliJ IDEA 平台中, 分为 UI 线程和后台线程, 这点和 Android 开发类似, 不同的是,

**读** 取操作可以在任何线程进行, 但在其他线程中读取需要使用 ***ApplicationManager.getApplication().runReadAction()*** 或者 ***ReadAction.run/compute*** 方法

**写** 操作只允许在 UI 线程进行, 必须使用 ***ApplicationManager.getApplication().runWriteAction()*** 或 ***WriteAction.run/compute*** 进行写操作

为了保证线程安全, 我们必须这样做

#### 什么是 PSI

PSI 是 Program Structure Interface 的缩写, 它定义了如何描述一种语言. 通过 AnActionEvent#getData(LangDataKeys.PSI_FILE) 获取当前文件的 PsiFile 对象.

每一种语言都有对应的 PsiFile 接口, 在插件开发模式下, 我们可以通过 **Tools => View PSI Structure** 查看一个文件的 PSI 结构, 他可以帮我们快速了解一种语言的 PSI 接口定义, 如果想开发解析某种语言的插件, 需要在项目中引入相应的 SDK.

Kotlin 类对应的 PSI 接口是 KtClass, 文件对应的是 KtFile

Java 类对应的 PSI 接口是 PsiClass, 文件对应的是 PsiJavaFile

一个源码文件的所有的元素都是 PsiElement 的子类, 包括 PsiFile, 比如在 Java 源码文件 PsiJavaFile 中 , 关键词 private, public 对应的 PsiElement 是 PsiKeyword. 通过PsiElement#acceptChild 方法可以遍历一个 element的所有子元素. 通过 PsiElement 的 add, delete, replace 等方法, 可以轻松的操作 PsiElement

创建一个用于 Java 的 PsiElement

	PsiElementFactory factory = JavaPsiFacade.getElementFactory(project);
	PsiField = factory.createFieldFromText("private String str = \"Hello\";", null);

创建一个用于 Kotlin 的 KtElement

	KtPsiFactory ktPsiFactory = KtPsiFactoryKt.KtPsiFactory(project);
	KtProperty ktProperty = ktPsiFactory.createProperty("private var str = \"Hello\"");

通过这两个工厂类可以创建所有的 PSI 元素, 当然我们也可以通过 new 实例化各种元素, 然后通过 add 关联在一起, 但这样相对比较麻烦.

#### 什么是 VFS

VFS 是 Virtual File System  的缩写, 它封装了大部分对活动文件的操作, 它提供了一个处理文件通用 API, 可以追踪文件变化

#### WriteCommandAction 操作 PSI

当我们使用 PsiElement#add(PsiElement e) 方法操作文件的时候需要用到这个类, WriteCommandAction#writeCommandAction(ThrowableRunnable t) 方法传入一个 Runnable.

#### 如何编写用户界面

我们可以选择 UI Designer, 或者自己手动敲. UI Designer 可以可视化编写界面, 直观, 在包目录上右键菜单 new 即可看到. 这个和 Swing 编程一毛一样. JetBrains 提供了它自己封装的一系列控件, 一般以 JB 开头, 比如 JBLabel, JBPanel, 有些特定的功能和统一的风格.

#### 技巧和注意事项

1. 插件开发, 我们需要使用 Community 版本的 IDEA, 否则无法调试源码

2. 如果没有发现 DevKit, 可能是该插件没有启用, 在 ***File > Settings > Plugins*** 中启用即可

3. 为了便于开发, 我们可以配置 IDEA 的源码, 在 https://github.com/JetBrains/intellij-community/ 仓库中下载与你 IDEA build 版本一支的源码, 然后添加到 ***ProjectStructure > SDKs > IntelliJ IDEA Community Edition IC xxx > Sourcepath***

4. 多个插件开发配置不同环境, 配置 ***SandBox ProjectStructure > SDKs > IntelliJ IDEA Community Edition IC xxx > Sandbox Home***

5. 导入插件项目不能直接 ***File> Open***, 而应该 ***File > New > Project From Existing Soruces...***

6. 在 ***Help > Edit Custom Properties*** 中 添加 ***idea.is.internal=true*** 并重启, 可以启用 ***Tools > Internal Actions***, 这里有许多好用的插件开发调试工具.

7. 工具 PSI Viewer ***Tools > View PSI Structure...*** 可以让我们快速了解到一个文件的 PSI 结构

## Action 的使用

Action 顾名思义就是动作, 用户可以通过按下一个快捷键或点击菜单选项触发.

### 定义

Action 定义了用户的一个动作, 快捷键, 我们创建一个 Action 需要一个类继承 AnAction, 并重写 actionPerformed(AnActionEvent anActionEvent) 方法, 之后在 plugin.xml 中注册该 Action.

基本上我们常用的数据上下文信息都可以在 anActionEvent 中获取, 例如光标: PlatformDataKeys.Carte, 获取当前语言 LangDataKeys.LANGUAGE.

例子, 定义一个 Action, 打印项目名, 路径, 及正在编辑的文件名

	public class MainAction extends AnAction {
		@Override
		public void actionPerformed(@NotNull AnActionEvent anActionEvent) {
			Project project = anActionEvent.getProject();
			PsiFile psiFile = anActionEvent.getData(LangDataKeys.PSI_FILE);
			System.out.println("Project Name:" + project.getName());
			System.out.println("Project Path" + project.getProjectFilePath());
			System.out.println("Editor File Name:" + psiFile.getName());
		}
	}

### 注册

在 plugin.xml 中注册该 action, 所有的 Action 都定义在 <actions></actions> 中

	<actions>
        <action id="your_id_usually_is_doaim_and_action_name" class="com.your_domain.MainAction"
                text="This is action name"
                description="This is description" keymap="$default">
            <add-to-group group-id="ToolsMenu" anchor="first"/> 
            <keyboard-shortcut first-keystroke="alt G" keymap="$default"/>
        </action>
	</actions>

group-id 定义了该 Action 出现的位置, 这里是在菜单 Tools 的第一个位置,  first-keystroke 为快捷键, 组合键用空格分开, 比如 "ctrl shift alt G".

我们在 Tools 第一个选项即可看到 "This is action name" 这个选项, 点击或按快捷键即可出发该 Action.

### 控制Action的隐藏显示

在一些情况, Action 在当前情况可能不可用, Action 是需要隐藏的, 比如, Generate=>toString 这个 Action 在编辑 xml 文件时就不适用, 需要隐藏, 重写 AnAction#update即可达到这个目的.

```java
public class ToStringAction extends AnAction {
    private static final LANG_XML = Language.findLanguageByID("XML");
    @Override
    public void update(@NotNull final AnActionEvent e) {
        Project project = e.getProject();
        PsiFile psiFile = e.getData(CommonDataKeys.PSI_FILE);
        
        e.getPresentation().setEnabledAndVisible(true);
        
        if(project == null || psiFile == null || !psiFile.getLanguage().is(LANG_XML)){
            e.getPresentation().setEnabledAndVisible(false);
        }
    }
}
```

以上代码可以实现没有打开 project, 没有打开文件或 语言不是 xml 时隐藏 ToStringAction.

## Editor

Editor 接口定义了对当前编辑器的一系列读写操作接口.

获取 Editor

	@Override
	public void actionPerformed(@NotNull AnActionEvent anActionEvent) {
		Editor editor = anActionEvent.getData(PlatformDataKeys.EDITOR);
	}

获取当前选择的文本

	SelectionModel selection = editor.getSelectionModel()
	if(selection != null){
		String text = selection.getSelectedText(true);
	}

Editor 可以获取一下8种 Model

- CaretModel 光标相关的 Model
- FoldingModel 折叠段落 Model
- IndentsModel	缩进 Model
- ScrollingModel 滚动 Model
- SoftWrapModel 自动换行 Model
- MarkupModel 标记,高亮 Model
- InlayModel 嵌套 Model
- SelectionModel 选择 Model

在获取相关 Model 时需要检查是否为空, 比如没有光标的时候, getCarteModel 将返回空. 针对我们要进行的不同操作获取不同的 Model.

## 组件

### ToolWindow

ToolWindow 就是底部 Logcat, Event Log 依附在左右两侧或底部的窗口, 可以最小化成一个按钮, 或展开, 改变大小和位置关闭.
在菜单栏中 View => Tool Window 列表中可以看到当前所有的 ToolWindow.

定义一个 ToolWindow, 显示当前项目名, 包上点击右键 new => Swing Ui Designer => GUI Form => TestToolWindow

点击 TestToolWindow.form 编辑界面, 添加一个 JLabel, 然后编辑 TestToolWindow, 让他实现 ToolWindowFactory 接口.

	public class TestToolWindow implements ToolWindowFactory {

		private JPanel rootPanel;
		private JLabel label1;

		public JPanel getContent() {
			return rootPanel;
		}

		@Override
		public void createToolWindowContent(@NotNull Project project, @NotNull ToolWindow toolWindow) {
			ContentFactory contentFactory = ContentFactory.SERVICE.getInstance();
			Content content = contentFactory.createContent(getContent(), "TestToolWindow", false);
			toolWindow.getContentManager().addContent(content);
		}
	}

在 plugin.xml 中注册, ToolWindow 需要放在 extensions 标签中.

    <extensions defaultExtensionNs="com.intellij">
        <toolWindow id="TestToolWindow"
                    canCloseContents="false"
                    factoryClass="com.your_domain.TestToolWindow"
                    anchor="bottom"/>
    </extensions>

其中, id 是 ToolWindow 的标题, canCloseContents 设置是否可以关闭, factoryClass 就是实现了 ToolWindowFactory 的该 ToolWindow 的工厂类. anchor 为显示位置

在 Action 中添加以下代码, 触发该 Action, ToolWindow 就弹出了并显示了项目的名称.

	public void actionPerformed(@NotNull AnActionEvent anActionEvent) {
	
		ToolWindow toolWindow = ToolWindowManager.getInstance(project).getToolWindow("TestToolWindow");
		toolWindow.show(new Runnable() {
			@Override
			public void run() {}
		});
		JTextField field = (JLabel) toolWindow.getContentManager()
				.getContent(0).getComponent().getComponent(0);
		if (field!=null){
			field.setText(project.getName());
		}
	}

### Dialog

IntelliJ SDK 中有一个 DialogWrap, 用这个可以与 IDEA 保持一致风格, 但是用这个就无法使用 GUI Designer 了. 它的使用方法与 Swing 中的 Dialog 差别不大.

一般情况, 我们开发的 plugin 都需要一个或若干个 Dialog.

### 持久化

PropertiesComponent 提供了数据持久化的接口, 他是一个单例, 通过 getInstance() 方法我们可以获取一个 Application 级的持久化实例, 他在所有的 Project 中都生效, 而 使用 PropertiesComponent.getInstance(Porject) 则只针对当前 Project 生效.

(完)